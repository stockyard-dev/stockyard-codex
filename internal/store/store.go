package store
import ("database/sql";"encoding/json";"fmt";"os";"path/filepath";"strings";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Snippet struct{ID string `json:"id"`;Title string `json:"title"`;Code string `json:"code"`;Language string `json:"language,omitempty"`;Description string `json:"description,omitempty"`;Tags []string `json:"tags"`;Public bool `json:"public"`;Favorite bool `json:"favorite"`;CreatedAt string `json:"created_at"`;UpdatedAt string `json:"updated_at"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"codex.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS snippets(id TEXT PRIMARY KEY,title TEXT NOT NULL,code TEXT DEFAULT '',language TEXT DEFAULT '',description TEXT DEFAULT '',tags_json TEXT DEFAULT '[]',public INTEGER DEFAULT 0,favorite INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')),updated_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(s *Snippet)error{s.ID=genID();s.CreatedAt=now();s.UpdatedAt=s.CreatedAt;if s.Tags==nil{s.Tags=[]string{}}
tj,_:=json.Marshal(s.Tags);pub:=0;if s.Public{pub=1};fav:=0;if s.Favorite{fav=1}
_,err:=d.db.Exec(`INSERT INTO snippets VALUES(?,?,?,?,?,?,?,?,?,?)`,s.ID,s.Title,s.Code,s.Language,s.Description,string(tj),pub,fav,s.CreatedAt,s.UpdatedAt);return err}
func(d *DB)scan(sc interface{Scan(...any)error})*Snippet{var s Snippet;var tj string;var pub,fav int
if sc.Scan(&s.ID,&s.Title,&s.Code,&s.Language,&s.Description,&tj,&pub,&fav,&s.CreatedAt,&s.UpdatedAt)!=nil{return nil}
json.Unmarshal([]byte(tj),&s.Tags);if s.Tags==nil{s.Tags=[]string{}};s.Public=pub==1;s.Favorite=fav==1;return &s}
func(d *DB)Get(id string)*Snippet{return d.scan(d.db.QueryRow(`SELECT * FROM snippets WHERE id=?`,id))}
func(d *DB)List(lang,tag string,favOnly bool)[]Snippet{where:=[]string{"1=1"};args:=[]any{}
if lang!=""{where=append(where,"language=?");args=append(args,lang)}
if tag!=""{where=append(where,`tags_json LIKE ?`);args=append(args,`%"`+tag+`"%`)}
if favOnly{where=append(where,"favorite=1")}
rows,_:=d.db.Query(`SELECT * FROM snippets WHERE `+strings.Join(where," AND ")+` ORDER BY favorite DESC,updated_at DESC`,args...);if rows==nil{return nil};defer rows.Close()
var o []Snippet;for rows.Next(){if s:=d.scan(rows);s!=nil{o=append(o,*s)}};return o}
func(d *DB)Update(id string,s *Snippet)error{tj,_:=json.Marshal(s.Tags);pub:=0;if s.Public{pub=1};fav:=0;if s.Favorite{fav=1}
_,err:=d.db.Exec(`UPDATE snippets SET title=?,code=?,language=?,description=?,tags_json=?,public=?,favorite=?,updated_at=? WHERE id=?`,s.Title,s.Code,s.Language,s.Description,string(tj),pub,fav,now(),id);return err}
func(d *DB)ToggleFavorite(id string)error{_,err:=d.db.Exec(`UPDATE snippets SET favorite=1-favorite,updated_at=? WHERE id=?`,now(),id);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM snippets WHERE id=?`,id);return err}
func(d *DB)Search(q string)[]Snippet{s:="%"+q+"%";rows,_:=d.db.Query(`SELECT * FROM snippets WHERE title LIKE ? OR code LIKE ? OR description LIKE ? ORDER BY updated_at DESC`,s,s,s);if rows==nil{return nil};defer rows.Close()
var o []Snippet;for rows.Next(){if sn:=d.scan(rows);sn!=nil{o=append(o,*sn)}};return o}
func(d *DB)Languages()[]string{rows,_:=d.db.Query(`SELECT DISTINCT language FROM snippets WHERE language!='' ORDER BY language`);if rows==nil{return nil};defer rows.Close();var o []string;for rows.Next(){var l string;rows.Scan(&l);o=append(o,l)};return o}
func(d *DB)AllTags()[]string{rows,_:=d.db.Query(`SELECT DISTINCT tags_json FROM snippets WHERE tags_json!='[]'`);if rows==nil{return nil};defer rows.Close()
seen:=map[string]bool{};for rows.Next(){var j string;rows.Scan(&j);var tags []string;json.Unmarshal([]byte(j),&tags);for _,t:=range tags{seen[t]=true}}
var o []string;for t:=range seen{o=append(o,t)};return o}
type Stats struct{Snippets int `json:"snippets"`;Languages int `json:"languages"`;Favorites int `json:"favorites"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM snippets`).Scan(&s.Snippets);s.Languages=len(d.Languages());d.db.QueryRow(`SELECT COUNT(*) FROM snippets WHERE favorite=1`).Scan(&s.Favorites);return s}
