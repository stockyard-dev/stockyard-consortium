package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Vendor struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ContactName string `json:"contact_name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Category string `json:"category"`
	ContractEnd string `json:"contract_end"`
	AnnualSpend int `json:"annual_spend"`
	Status string `json:"status"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"consortium.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS vendors(id TEXT PRIMARY KEY,name TEXT NOT NULL,contact_name TEXT DEFAULT '',email TEXT DEFAULT '',phone TEXT DEFAULT '',category TEXT DEFAULT '',contract_end TEXT DEFAULT '',annual_spend INTEGER DEFAULT 0,status TEXT DEFAULT 'active',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Vendor)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO vendors(id,name,contact_name,email,phone,category,contract_end,annual_spend,status,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.ContactName,e.Email,e.Phone,e.Category,e.ContractEnd,e.AnnualSpend,e.Status,e.Notes,e.CreatedAt);return err}
func(d *DB)Get(id string)*Vendor{var e Vendor;if d.db.QueryRow(`SELECT id,name,contact_name,email,phone,category,contract_end,annual_spend,status,notes,created_at FROM vendors WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.ContactName,&e.Email,&e.Phone,&e.Category,&e.ContractEnd,&e.AnnualSpend,&e.Status,&e.Notes,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Vendor{rows,_:=d.db.Query(`SELECT id,name,contact_name,email,phone,category,contract_end,annual_spend,status,notes,created_at FROM vendors ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Vendor;for rows.Next(){var e Vendor;rows.Scan(&e.ID,&e.Name,&e.ContactName,&e.Email,&e.Phone,&e.Category,&e.ContractEnd,&e.AnnualSpend,&e.Status,&e.Notes,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM vendors WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM vendors`).Scan(&n);return n}
