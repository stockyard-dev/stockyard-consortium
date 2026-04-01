package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-consortium/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){q:=r.URL.Query().Get("q");list,_:=s.db.List(q);if list==nil{list=[]store.Vendor{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var v store.Vendor;json.NewDecoder(r.Body).Decode(&v);if v.Name==""{writeError(w,400,"name required");return};s.db.Create(&v);writeJSON(w,201,v)}
func(s *Server)handleUpdate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var v store.Vendor;json.NewDecoder(r.Body).Decode(&v);v.ID=id;s.db.Update(&v);writeJSON(w,200,v)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
