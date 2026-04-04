package server
import "net/http"
func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html"); w.Write([]byte(dashHTML)) }
const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Consortium</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}.main{padding:1.5rem;max-width:960px;margin:0 auto}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.search:focus{outline:none;border-color:var(--leather)}.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}.item:hover{border-color:var(--leather)}.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}.item-title{font-size:.85rem;font-weight:700}.item-sub{font-size:.7rem;color:var(--cd);margin-top:.1rem}.item-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}.item-actions{display:flex;gap:.3rem;flex-shrink:0}.badge{font-size:.5rem;padding:.12rem .35rem;border:1px solid var(--bg3);color:var(--cm)}.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-sm{font-size:.55rem;padding:.2rem .4rem}.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.fr input:focus{outline:none;border-color:var(--leather)}.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> CONSORTIUM</h1><button class="btn btn-p" onclick="openForm()">+ Add Vendor</button></div>
<div class="main"><div class="stats" id="stats"></div><div class="toolbar"><input class="search" id="search" placeholder="Search vendors..." oninput="render()"></div><div id="list"></div></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/vendors').then(function(r){return r.json()});items=r.vendors||[];renderStats();render();}
function renderStats(){var total=items.length;var spend=items.reduce(function(s,v){return s+(v.annual_spend||0)},0);var cats={};items.forEach(function(v){if(v.category)cats[v.category]=true});
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Vendors</div></div><div class="st"><div class="st-v">$'+(spend/100).toLocaleString()+'</div><div class="st-l">Annual Spend</div></div><div class="st"><div class="st-v">'+Object.keys(cats).length+'</div><div class="st-l">Categories</div></div>';}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(v){return(v.name||'').toLowerCase().includes(q)||(v.contact_name||'').toLowerCase().includes(q)||(v.category||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No vendors.</div>';return;}
var h='';f.forEach(function(v){
h+='<div class="item"><div class="item-top"><div style="flex:1"><div class="item-title">'+esc(v.name)+'</div>';
if(v.contact_name)h+='<div class="item-sub">'+esc(v.contact_name)+(v.email?' &#183; '+esc(v.email):'')+'</div>';
h+='</div><div class="item-actions"><button class="btn btn-sm" onclick="openEdit(''+v.id+'')">Edit</button><button class="btn btn-sm" onclick="del(''+v.id+'')" style="color:var(--red)">&#10005;</button></div></div>';
h+='<div class="item-meta">';
if(v.category)h+='<span class="badge">'+esc(v.category)+'</span>';
if(v.annual_spend)h+='<span>$'+(v.annual_spend/100).toFixed(2)+'/yr</span>';
if(v.contract_end)h+='<span>Contract: '+v.contract_end+'</span>';
if(v.phone)h+='<span>'+esc(v.phone)+'</span>';
h+='</div></div>';});document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/vendors/'+id,{method:'DELETE'});load();}
function formHTML(v){var i=v||{name:'',contact_name:'',email:'',phone:'',category:'',contract_end:'',annual_spend:0};var isEdit=!!v;
var h='<h2>'+(isEdit?'EDIT':'ADD')+' VENDOR</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'"></div>';
h+='<div class="row2"><div class="fr"><label>Contact</label><input id="f-contact" value="'+esc(i.contact_name)+'"></div><div class="fr"><label>Email</label><input id="f-email" value="'+esc(i.email)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Phone</label><input id="f-phone" value="'+esc(i.phone)+'"></div><div class="fr"><label>Category</label><input id="f-cat" value="'+esc(i.category)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Contract End</label><input id="f-end" type="date" value="'+esc(i.contract_end)+'"></div><div class="fr"><label>Annual Spend ($)</label><input id="f-spend" type="number" step="0.01" value="'+((i.annual_spend||0)/100).toFixed(2)+'"></div></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var v=null;for(var j=0;j<items.length;j++){if(items[j].id===id){v=items[j];break;}}if(!v)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(v);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var name=document.getElementById('f-name').value.trim();if(!name){alert('Name required');return;}
var body={name:name,contact_name:document.getElementById('f-contact').value.trim(),email:document.getElementById('f-email').value.trim(),phone:document.getElementById('f-phone').value.trim(),category:document.getElementById('f-cat').value.trim(),contract_end:document.getElementById('f-end').value,annual_spend:Math.round(parseFloat(document.getElementById('f-spend').value||0)*100)};
if(editId){await fetch(A+'/vendors/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/vendors',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}closeModal();load();}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});load();
</script></body></html>`
