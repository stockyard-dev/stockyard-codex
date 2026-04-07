package server

import (
	"encoding/json"
	"github.com/stockyard-dev/stockyard-codex/internal/store"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	db      *store.DB
	mux     *http.ServeMux
	limits  Limits
	dataDir string
	pCfg    map[string]json.RawMessage
}

func New(db *store.DB, limits Limits, dataDir string) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), limits: limits, dataDir: dataDir}
	s.mux.HandleFunc("GET /api/snippets", s.list)
	s.mux.HandleFunc("POST /api/snippets", s.create)
	s.mux.HandleFunc("GET /api/snippets/{id}", s.get)
	s.mux.HandleFunc("PUT /api/snippets/{id}", s.update)
	s.mux.HandleFunc("DELETE /api/snippets/{id}", s.del)
	s.mux.HandleFunc("POST /api/snippets/{id}/favorite", s.toggleFav)
	s.mux.HandleFunc("GET /api/search", s.search)
	s.mux.HandleFunc("GET /api/languages", s.languages)
	s.mux.HandleFunc("GET /api/tags", s.tags)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{"tier": s.limits.Tier, "upgrade_url": "https://stockyard.dev/codex/"})
	})
	s.loadPersonalConfig()
	s.mux.HandleFunc("GET /api/config", s.configHandler)
	s.mux.HandleFunc("GET /api/extras/{resource}", s.listExtras)
	s.mux.HandleFunc("GET /api/extras/{resource}/{id}", s.getExtras)
	s.mux.HandleFunc("PUT /api/extras/{resource}/{id}", s.putExtras)
	return s
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }
func wj(w http.ResponseWriter, c int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	json.NewEncoder(w).Encode(v)
}
func we(w http.ResponseWriter, c int, m string) { wj(w, c, map[string]string{"error": m}) }
func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui", 302)
}
func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	wj(w, 200, map[string]any{"snippets": oe(s.db.List(q.Get("language"), q.Get("tag"), q.Get("favorites") == "true"))})
}
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	var sn store.Snippet
	json.NewDecoder(r.Body).Decode(&sn)
	if sn.Title == "" {
		we(w, 400, "title required")
		return
	}
	s.db.Create(&sn)
	wj(w, 201, s.db.Get(sn.ID))
}
func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	sn := s.db.Get(r.PathValue("id"))
	if sn == nil {
		we(w, 404, "not found")
		return
	}
	wj(w, 200, sn)
}
func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ex := s.db.Get(id)
	if ex == nil {
		we(w, 404, "not found")
		return
	}
	var sn store.Snippet
	json.NewDecoder(r.Body).Decode(&sn)
	if sn.Title == "" {
		sn.Title = ex.Title
	}
	if sn.Tags == nil {
		sn.Tags = ex.Tags
	}
	s.db.Update(id, &sn)
	wj(w, 200, s.db.Get(id))
}
func (s *Server) del(w http.ResponseWriter, r *http.Request) {
	s.db.Delete(r.PathValue("id"))
	wj(w, 200, map[string]string{"deleted": "ok"})
}
func (s *Server) toggleFav(w http.ResponseWriter, r *http.Request) {
	s.db.ToggleFavorite(r.PathValue("id"))
	wj(w, 200, s.db.Get(r.PathValue("id")))
}
func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, map[string]any{"snippets": oe(s.db.Search(r.URL.Query().Get("q")))})
}
func (s *Server) languages(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, map[string]any{"languages": oe(s.db.Languages())})
}
func (s *Server) tags(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, map[string]any{"tags": oe(s.db.AllTags())})
}
func (s *Server) stats(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) }
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats()
	wj(w, 200, map[string]any{"status": "ok", "service": "codex", "snippets": st.Snippets})
}
func oe[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
func init() { log.SetFlags(log.LstdFlags | log.Lshortfile) }

// ─── personalization (auto-added) ──────────────────────────────────

func (s *Server) loadPersonalConfig() {
	path := filepath.Join(s.dataDir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("%s: warning: could not parse config.json: %v", "codex", err)
		return
	}
	s.pCfg = cfg
	log.Printf("%s: loaded personalization from %s", "codex", path)
}

func (s *Server) configHandler(w http.ResponseWriter, r *http.Request) {
	if s.pCfg == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.pCfg)
}

func (s *Server) listExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	all := s.db.AllExtras(resource)
	out := make(map[string]json.RawMessage, len(all))
	for id, data := range all {
		out[id] = json.RawMessage(data)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (s *Server) getExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	data := s.db.GetExtras(resource, id)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func (s *Server) putExtras(w http.ResponseWriter, r *http.Request) {
	resource := r.PathValue("resource")
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"read body"}`, 400)
		return
	}
	var probe map[string]any
	if err := json.Unmarshal(body, &probe); err != nil {
		http.Error(w, `{"error":"invalid json"}`, 400)
		return
	}
	if err := s.db.SetExtras(resource, id, string(body)); err != nil {
		http.Error(w, `{"error":"save failed"}`, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":"saved"}`))
}
