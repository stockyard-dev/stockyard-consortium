package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Vendor is a tracked supplier or service provider. AnnualSpend is stored
// as integer cents to avoid floating-point money bugs (matches the steward
// pattern). The display layer converts to dollars for the UI.
type Vendor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ContactName string `json:"contact_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Category    string `json:"category"`
	ContractEnd string `json:"contract_end"`
	AnnualSpend int    `json:"annual_spend"` // cents
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "consortium.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS vendors(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		contact_name TEXT DEFAULT '',
		email TEXT DEFAULT '',
		phone TEXT DEFAULT '',
		category TEXT DEFAULT '',
		contract_end TEXT DEFAULT '',
		annual_spend INTEGER DEFAULT 0,
		status TEXT DEFAULT 'active',
		notes TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_vendors_status ON vendors(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_vendors_category ON vendors(category)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Vendor) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Status == "" {
		e.Status = "active"
	}
	_, err := d.db.Exec(
		`INSERT INTO vendors(id, name, contact_name, email, phone, category, contract_end, annual_spend, status, notes, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.ContactName, e.Email, e.Phone, e.Category, e.ContractEnd, e.AnnualSpend, e.Status, e.Notes, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Vendor {
	var e Vendor
	err := d.db.QueryRow(
		`SELECT id, name, contact_name, email, phone, category, contract_end, annual_spend, status, notes, created_at
		 FROM vendors WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Name, &e.ContactName, &e.Email, &e.Phone, &e.Category, &e.ContractEnd, &e.AnnualSpend, &e.Status, &e.Notes, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Vendor {
	rows, _ := d.db.Query(
		`SELECT id, name, contact_name, email, phone, category, contract_end, annual_spend, status, notes, created_at
		 FROM vendors ORDER BY name ASC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Vendor
	for rows.Next() {
		var e Vendor
		rows.Scan(&e.ID, &e.Name, &e.ContactName, &e.Email, &e.Phone, &e.Category, &e.ContractEnd, &e.AnnualSpend, &e.Status, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Vendor) error {
	_, err := d.db.Exec(
		`UPDATE vendors SET name=?, contact_name=?, email=?, phone=?, category=?, contract_end=?, annual_spend=?, status=?, notes=?
		 WHERE id=?`,
		e.Name, e.ContactName, e.Email, e.Phone, e.Category, e.ContractEnd, e.AnnualSpend, e.Status, e.Notes, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM vendors WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM vendors`).Scan(&n)
	return n
}

// Search filters vendors by query (name/contact/email) and optional
// category and status filters.
func (d *DB) Search(q string, filters map[string]string) []Vendor {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ? OR contact_name LIKE ? OR email LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s)
	}
	if v, ok := filters["category"]; ok && v != "" {
		where += " AND category=?"
		args = append(args, v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, name, contact_name, email, phone, category, contract_end, annual_spend, status, notes, created_at
		 FROM vendors WHERE `+where+`
		 ORDER BY name ASC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Vendor
	for rows.Next() {
		var e Vendor
		rows.Scan(&e.ID, &e.Name, &e.ContactName, &e.Email, &e.Phone, &e.Category, &e.ContractEnd, &e.AnnualSpend, &e.Status, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns aggregate metrics for the dashboard. Includes total vendor
// count, total annual spend (cents), counts by status and category, and
// the number of contracts expiring within 30 days.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":         d.Count(),
		"total_spend":   0,
		"by_status":     map[string]int{},
		"by_category":   map[string]int{},
		"expiring_soon": 0,
	}

	var totalSpend int
	d.db.QueryRow(`SELECT COALESCE(SUM(annual_spend), 0) FROM vendors WHERE status='active'`).Scan(&totalSpend)
	m["total_spend"] = totalSpend

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM vendors GROUP BY status`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}

	if rows, _ := d.db.Query(`SELECT category, COUNT(*) FROM vendors WHERE category != '' GROUP BY category`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_category"] = by
	}

	// Contracts expiring within 30 days (and not already past)
	today := time.Now().Format("2006-01-02")
	thirtyDays := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	var expiring int
	d.db.QueryRow(
		`SELECT COUNT(*) FROM vendors WHERE contract_end != '' AND contract_end >= ? AND contract_end <= ? AND status='active'`,
		today, thirtyDays,
	).Scan(&expiring)
	m["expiring_soon"] = expiring

	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
