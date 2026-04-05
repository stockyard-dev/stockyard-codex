package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-codex/internal/server";"github.com/stockyard-dev/stockyard-codex/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8650"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./codex-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("codex: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Codex — Self-hosted code snippet manager\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n  Questions? hello@stockyard.dev\n\n",port,port,dataDir)
log.Printf("codex: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
