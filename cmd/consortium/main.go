package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-consortium/internal/server";"github.com/stockyard-dev/stockyard-consortium/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9850"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./consortium-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("consortium: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Consortium — vendor and supplier tracker\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("consortium: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
