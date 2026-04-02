package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-codex/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux}
func New(db *store.DB)*Server{s:=&Server{db:db,mux:http.NewServeMux()}
s.mux.HandleFunc("GET /api/snippets",s.list);s.mux.HandleFunc("POST /api/snippets",s.create);s.mux.HandleFunc("GET /api/snippets/{id}",s.get);s.mux.HandleFunc("PUT /api/snippets/{id}",s.update);s.mux.HandleFunc("DELETE /api/snippets/{id}",s.del)
s.mux.HandleFunc("POST /api/snippets/{id}/favorite",s.toggleFav)
s.mux.HandleFunc("GET /api/search",s.search);s.mux.HandleFunc("GET /api/languages",s.languages);s.mux.HandleFunc("GET /api/tags",s.tags)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){q:=r.URL.Query();wj(w,200,map[string]any{"snippets":oe(s.db.List(q.Get("language"),q.Get("tag"),q.Get("favorites")=="true"))})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var sn store.Snippet;json.NewDecoder(r.Body).Decode(&sn);if sn.Title==""{we(w,400,"title required");return};s.db.Create(&sn);wj(w,201,s.db.Get(sn.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){sn:=s.db.Get(r.PathValue("id"));if sn==nil{we(w,404,"not found");return};wj(w,200,sn)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.Get(id);if ex==nil{we(w,404,"not found");return};var sn store.Snippet;json.NewDecoder(r.Body).Decode(&sn);if sn.Title==""{sn.Title=ex.Title};if sn.Tags==nil{sn.Tags=ex.Tags};s.db.Update(id,&sn);wj(w,200,s.db.Get(id))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)toggleFav(w http.ResponseWriter,r *http.Request){s.db.ToggleFavorite(r.PathValue("id"));wj(w,200,s.db.Get(r.PathValue("id")))}
func(s *Server)search(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"snippets":oe(s.db.Search(r.URL.Query().Get("q")))})}
func(s *Server)languages(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"languages":oe(s.db.Languages())})}
func(s *Server)tags(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tags":oe(s.db.AllTags())})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"codex","snippets":st.Snippets})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
