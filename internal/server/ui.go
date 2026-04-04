package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Codex</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}
.main{padding:1.5rem;max-width:900px;margin:0 auto}
.search{width:100%;padding:.5rem .8rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.78rem;margin-bottom:1rem}
.lang-bar{display:flex;gap:.3rem;margin-bottom:1rem;flex-wrap:wrap}
.lang-btn{font-size:.6rem;padding:.2rem .5rem;border:1px solid var(--bg3);background:var(--bg);color:var(--cm);cursor:pointer}.lang-btn:hover{border-color:var(--leather)}.lang-btn.active{border-color:var(--rust);color:var(--rust)}
.snip{background:var(--bg2);border:1px solid var(--bg3);margin-bottom:.8rem}
.snip-hdr{padding:.6rem .8rem;display:flex;justify-content:space-between;align-items:center;border-bottom:1px solid var(--bg3)}
.snip-title{font-size:.82rem;color:var(--cream)}.snip-lang{font-size:.55rem;padding:.1rem .3rem;background:var(--bg3);color:var(--gold)}
.snip-code{padding:.8rem;background:#0d0b09;font-size:.72rem;color:var(--cd);overflow-x:auto;white-space:pre;max-height:200px;cursor:pointer;position:relative}
.snip-code:hover::after{content:'click to copy';position:absolute;top:.3rem;right:.3rem;font-size:.5rem;color:var(--cm)}
.snip-meta{padding:.4rem .8rem;font-size:.55rem;color:var(--cm);display:flex;gap:.6rem}
.tag{font-size:.5rem;padding:.05rem .25rem;background:var(--bg3);color:var(--cm)}
.btn{font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:500px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr textarea{min-height:150px;white-space:pre}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>CODEX</h1><button class="btn btn-p" onclick="openForm()">+ New Snippet</button></div>
<div class="main">
<input class="search" id="search" placeholder="Search snippets..." oninput="render()">
<div class="lang-bar" id="langs"></div>
<div id="snippets"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let snippets=[],filterLang='';
async function load(){const r=await fetch(A+'/snippets').then(r=>r.json());snippets=r.snippets||[];
const langs=[...new Set(snippets.map(s=>s.language).filter(l=>l))];
let lh='<button class="lang-btn'+(filterLang===''?' active':'')+'" onclick="setLang(\'\')">All ('+snippets.length+')</button>';
langs.forEach(l=>{lh+='<button class="lang-btn'+(filterLang===l?' active':'')+'" onclick="setLang(\''+l+'\')">'+esc(l)+'</button>';});
document.getElementById('langs').innerHTML=lh;render();}
function setLang(l){filterLang=l;render();}
function render(){const q=(document.getElementById('search').value||'').toLowerCase();
let filtered=snippets.filter(s=>{if(filterLang&&s.language!==filterLang)return false;if(q&&!(s.title+s.code+s.description+s.language).toLowerCase().includes(q))return false;return true;});
if(!filtered.length){document.getElementById('snippets').innerHTML='<div class="empty">No snippets.</div>';return;}
let h='';filtered.forEach(s=>{
h+='<div class="snip"><div class="snip-hdr"><div class="snip-title">'+(s.favorite?'★ ':'')+esc(s.title)+'</div><div style="display:flex;gap:.3rem;align-items:center">';
if(s.language)h+='<span class="snip-lang">'+esc(s.language)+'</span>';
h+='<button class="btn" onclick="del(\''+s.id+'\')" style="font-size:.5rem;color:var(--cm)">✕</button></div></div>';
h+='<div class="snip-code" onclick="navigator.clipboard.writeText(this.textContent)">'+esc(s.code)+'</div>';
h+='<div class="snip-meta">';if(s.description)h+='<span>'+esc(s.description)+'</span>';
const tags=JSON.parse(s.tags_json||'[]');tags.forEach(t=>{h+='<span class="tag">'+esc(t)+'</span>';});
h+='<span>'+ft(s.created_at)+'</span></div></div>';});
document.getElementById('snippets').innerHTML=h;}
async function del(id){if(confirm('Delete?')){await fetch(A+'/snippets/'+id,{method:'DELETE'});load();}}
function openForm(){document.getElementById('mdl').innerHTML='<h2>New Snippet</h2><div class="fr"><label>Title</label><input id="f-t" placeholder="e.g. Retry with exponential backoff"></div><div class="fr"><label>Language</label><input id="f-l" placeholder="go, python, javascript, sql"></div><div class="fr"><label>Code</label><textarea id="f-c" placeholder="paste your code here"></textarea></div><div class="fr"><label>Description</label><input id="f-d"></div><div class="fr"><label>Tags (JSON array)</label><input id="f-tg" value="[]"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Save</button></div>';document.getElementById('mbg').classList.add('open');}
async function sub(){await fetch(A+'/snippets',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({title:document.getElementById('f-t').value,language:document.getElementById('f-l').value,code:document.getElementById('f-c').value,description:document.getElementById('f-d').value,tags_json:document.getElementById('f-tg').value})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function ft(t){if(!t)return'';return new Date(t).toLocaleDateString();}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
