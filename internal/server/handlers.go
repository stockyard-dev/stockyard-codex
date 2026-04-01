package server
import("encoding/json";"fmt";"net/http";"regexp";"strconv";"strings";"github.com/stockyard-dev/stockyard-codex/internal/store")
var slugRe=regexp.MustCompile(`[^a-z0-9-]`)
func toSlug(s string)string{return slugRe.ReplaceAllString(strings.ToLower(strings.ReplaceAll(s," ","-")),"")  }
func(s *Server)handleListPages(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListPages();if list==nil{list=[]store.Page{}};writeJSON(w,200,list)}
func(s *Server)handleGetPage(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);p,_:=s.db.GetPage(id);if p==nil{writeError(w,404,"not found");return};writeJSON(w,200,p)}
func(s *Server)handleCreatePage(w http.ResponseWriter,r *http.Request){
    if !s.limits.IsPro(){n,_:=s.db.CountPages();if n>=20{writeError(w,403,"free tier: 20 pages max");return}}
    var p store.Page;json.NewDecoder(r.Body).Decode(&p)
    if p.Title==""{writeError(w,400,"title required");return}
    if p.Slug==""{p.Slug=toSlug(p.Title)}
    if err:=s.db.CreatePage(&p);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,201,p)}
func(s *Server)handleUpdatePage(w http.ResponseWriter,r *http.Request){
    id,_:=strconv.ParseInt(r.PathValue("id"),10,64)
    existing,_:=s.db.GetPage(id);if existing==nil{writeError(w,404,"not found");return}
    json.NewDecoder(r.Body).Decode(existing);existing.ID=id
    if existing.Slug==""{existing.Slug=toSlug(existing.Title)}
    if err:=s.db.UpdatePage(existing);err!=nil{writeError(w,500,err.Error());return}
    writeJSON(w,200,existing)}
func(s *Server)handleDeletePage(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeletePage(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListVersions(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListVersions(id);if list==nil{list=[]store.PageVersion{}};writeJSON(w,200,list)}
func(s *Server)handleSearch(w http.ResponseWriter,r *http.Request){q:=r.URL.Query().Get("q");if q==""{writeJSON(w,200,[]store.Page{});return};list,err:=s.db.SearchPages(q);if err!=nil{list,_=s.db.ListPages()};if list==nil{list=[]store.Page{}};writeJSON(w,200,list)}
func(s *Server)handleWikiPage(w http.ResponseWriter,r *http.Request){
    slug:=r.PathValue("slug");p,_:=s.db.GetPageBySlug(slug);if p==nil{http.NotFound(w,r);return}
    w.Header().Set("Content-Type","text/html")
    content:=strings.ReplaceAll(p.Content,"\n","<br>")
    fmt.Fprintf(w,`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>%s — Codex</title><style>body{background:#1a1410;color:#e8d5b0;font-family:serif;max-width:800px;margin:2rem auto;padding:1rem;line-height:1.7}h1{color:#c4622d;font-size:2rem}a{color:#8b5e3c}.meta{color:#7a6550;font-size:0.85rem;margin-bottom:2rem}.tags span{background:#241c15;border:1px solid #3d2e1e;border-radius:3px;padding:0.1rem 0.4rem;font-size:0.78rem;margin-right:0.3rem}</style></head><body><h1>%s</h1><div class="meta">Updated %s &bull; <a href="/">← Back to wiki</a></div>`,p.Title,p.Title,p.UpdatedAt.Format("Jan 2, 2006"))
    if p.Tags!=""{fmt.Fprintf(w,`<div class="tags">Tags: `)
        for _,t:=range strings.Split(p.Tags,","){fmt.Fprintf(w,`<span>%s</span>`,strings.TrimSpace(t))};fmt.Fprintf(w,`</div><br>`)}
    fmt.Fprintf(w,`<div>%s</div></body></html>`,content)}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){p,_:=s.db.CountPages();writeJSON(w,200,map[string]interface{}{"pages":p,"versions":0})}
