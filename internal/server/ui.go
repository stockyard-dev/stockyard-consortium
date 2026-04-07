package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Consortium</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1000px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-v.warn{color:var(--orange)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .15s}
.item:hover{border-color:var(--leather)}
.item.expiring{border-left:3px solid var(--orange)}
.item.inactive{opacity:.6}
.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.item-title{font-size:.85rem;font-weight:700;color:var(--cream)}
.item-sub{font-size:.65rem;color:var(--cd);margin-top:.15rem}
.item-meta{font-size:.55rem;color:var(--cm);margin-top:.4rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.item-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.item-extra-row{display:flex;gap:.4rem}
.item-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.item-extra-val{color:var(--cream)}
.item-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.active{border-color:var(--green);color:var(--green)}
.badge.inactive,.badge.terminated{border-color:var(--cm);color:var(--cm)}
.badge.on_hold,.badge.pending{border-color:var(--orange);color:var(--orange)}
.badge.cat{border-color:var(--leather);color:var(--leather)}
.spend{color:var(--gold);font-weight:700}
.expiring-tag{color:var(--orange);font-weight:700}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:520px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto;color:var(--red);border-color:#3a1a1a}
.acts .btn-del:hover{border-color:var(--red);color:var(--red)}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> CONSORTIUM</h1>
<button class="btn btn-p" onclick="openForm()">+ Add Vendor</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search vendors..." oninput="debouncedRender()">
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Statuses</option>
<option value="active">Active</option>
<option value="on_hold">On Hold</option>
<option value="terminated">Terminated</option>
</select>
<select class="filter-sel" id="category-filter" onchange="render()">
<option value="">All Categories</option>
</select>
</div>
<div id="list"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='vendors';

// Field defs for vendor form. Custom fields injected from /api/config get
// isCustom=true and persist to the extras table.
var fields=[
{name:'name',label:'Company Name',type:'text',required:true},
{name:'contact_name',label:'Contact',type:'text'},
{name:'email',label:'Email',type:'email'},
{name:'phone',label:'Phone',type:'text'},
{name:'category',label:'Category',type:'select_or_text',options:[]},
{name:'status',label:'Status',type:'select',options:['active','on_hold','terminated']},
{name:'contract_end',label:'Contract End',type:'date'},
{name:'annual_spend',label:'Annual Spend ($)',type:'money'},
{name:'notes',label:'Notes',type:'textarea'}
];

var vendors=[],vendorExtras={},editId=null,searchTimer=null;

// ─── Helpers ──────────────────────────────────────────────────────

// Money helpers — annual_spend is stored as cents
function fmtMoney(cents){
if(!cents)return'$0';
var dollars=cents/100;
if(dollars>=1000000)return'$'+(dollars/1000000).toFixed(1)+'M';
if(dollars>=1000)return'$'+(dollars/1000).toFixed(1)+'k';
return'$'+dollars.toFixed(0);
}
function fmtMoneyFull(cents){
if(!cents)return'$0';
return'$'+(cents/100).toLocaleString('en-US',{minimumFractionDigits:0,maximumFractionDigits:0});
}
function parseMoney(str){
if(!str)return 0;
var n=parseFloat(String(str).replace(/[^\d.-]/g,''));
if(isNaN(n))return 0;
return Math.round(n*100);
}

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s+'T12:00:00');
if(isNaN(d.getTime()))return s;
return d.toLocaleDateString('en-US',{year:'numeric',month:'short',day:'numeric'});
}catch(e){return s}
}

function isExpiringSoon(dateStr,status){
if(!dateStr||status!=='active')return false;
try{
var end=new Date(dateStr+'T12:00:00');
var now=new Date();
var diff=(end-now)/(1000*60*60*24);
return diff>=0&&diff<=30;
}catch(e){return false}
}

function isExpired(dateStr){
if(!dateStr)return false;
try{
var end=new Date(dateStr+'T12:00:00');
return end<new Date();
}catch(e){return false}
}

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/vendors').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
vendors=resps[0].vendors||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
vendorExtras=ex||{};
vendors.forEach(function(v){
var x=vendorExtras[v.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(v[k]===undefined)v[k]=x[k]});
});
}catch(e){vendorExtras={}}

populateCategoryFilter();
}catch(e){
console.error('load failed',e);
vendors=[];
}
render();
}

function populateCategoryFilter(){
var sel=document.getElementById('category-filter');
if(!sel)return;
var current=sel.value;
var seen={};
var cats=[];
vendors.forEach(function(v){
if(v.category&&!seen[v.category]){seen[v.category]=true;cats.push(v.category)}
});
cats.sort();
sel.innerHTML='<option value="">All Categories</option>'+cats.map(function(c){return'<option value="'+esc(c)+'"'+(c===current?' selected':'')+'>'+esc(c)+'</option>'}).join('');
}

function renderStats(s){
var total=s.total||0;
var spend=s.total_spend||0;
var expiring=s.expiring_soon||0;
var byCat=s.by_category||{};
var catCount=Object.keys(byCat).length;
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Vendors</div></div>'+
'<div class="st"><div class="st-v">'+fmtMoney(spend)+'</div><div class="st-l">Annual Spend</div></div>'+
'<div class="st"><div class="st-v '+(expiring>0?'warn':'')+'">'+expiring+'</div><div class="st-l">Expiring 30d</div></div>'+
'<div class="st"><div class="st-v">'+catCount+'</div><div class="st-l">Categories</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var cf=document.getElementById('category-filter').value;

var f=vendors;
if(q)f=f.filter(function(v){
return(v.name||'').toLowerCase().includes(q)||
(v.contact_name||'').toLowerCase().includes(q)||
(v.email||'').toLowerCase().includes(q)||
(v.category||'').toLowerCase().includes(q);
});
if(sf)f=f.filter(function(v){return v.status===sf});
if(cf)f=f.filter(function(v){return v.category===cf});

if(!f.length){
var msg=window._emptyMsg||'No vendors. Add your first one.';
document.getElementById('list').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(v){h+=itemHTML(v)});
document.getElementById('list').innerHTML=h;
}

function itemHTML(v){
var expiring=isExpiringSoon(v.contract_end,v.status);
var expired=isExpired(v.contract_end)&&v.status==='active';
var cls='item';
if(expiring||expired)cls+=' expiring';
if(v.status==='terminated'||v.status==='inactive')cls+=' inactive';

var h='<div class="'+cls+'">';
h+='<div class="item-top"><div style="flex:1;min-width:0">';
h+='<div class="item-title">'+esc(v.name)+'</div>';
if(v.contact_name||v.email){
h+='<div class="item-sub">';
if(v.contact_name)h+=esc(v.contact_name);
if(v.email)h+=(v.contact_name?' &middot; ':'')+esc(v.email);
h+='</div>';
}
h+='</div><div class="item-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(\''+esc(v.id)+'\')">Edit</button>';
h+='</div></div>';

h+='<div class="item-meta">';
if(v.status)h+='<span class="badge '+esc(v.status)+'">'+esc(v.status.replace(/_/g,' '))+'</span>';
if(v.category)h+='<span class="badge cat">'+esc(v.category)+'</span>';
if(v.annual_spend)h+='<span class="spend">'+fmtMoneyFull(v.annual_spend)+'/yr</span>';
if(v.phone)h+='<span>'+esc(v.phone)+'</span>';
if(v.contract_end){
var label='Contract: '+fmtDate(v.contract_end);
if(expired)h+='<span class="expiring-tag">EXPIRED ('+fmtDate(v.contract_end)+')</span>';
else if(expiring)h+='<span class="expiring-tag">'+label+' (soon)</span>';
else h+='<span>'+esc(label)+'</span>';
}
h+='</div>';

// Custom fields display
var customRows='';
fields.forEach(function(f){
if(!f.isCustom)return;
var val=v[f.name];
if(val===undefined||val===null||val==='')return;
customRows+='<div class="item-extra-row">';
customRows+='<span class="item-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="item-extra-val">'+esc(String(val))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="item-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

// ─── Modal: vendor form ───────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
var disp=String(o).charAt(0).toUpperCase()+String(o).slice(1).replace(/_/g,' ');
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(disp)+'</option>';
});
h+='</select>';
}else if(f.type==='select_or_text'){
// Datalist input — text input with autocomplete from existing values
h+='<input list="dl-'+f.name+'" type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
h+='<datalist id="dl-'+f.name+'">';
var opts=(f.options||[]).slice();
// Inject existing categories from vendors
vendors.forEach(function(vd){
if(vd.category&&opts.indexOf(vd.category)===-1)opts.push(vd.category);
});
opts.forEach(function(o){h+='<option value="'+esc(String(o))+'">'});
h+='</datalist>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="2"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='money'){
var dollars=v?(v/100).toFixed(2):'';
h+='<input type="text" id="f-'+f.name+'" value="'+esc(dollars)+'" placeholder="0.00">';
}else if(f.type==='number'||f.type==='integer'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(vendor){
var v=vendor||{};
var isEdit=!!vendor;
var h='<h2>'+(isEdit?'EDIT VENDOR':'NEW VENDOR')+'</h2>';

h+=fieldHTML(fieldByName('name'),v.name);
h+='<div class="row2">'+fieldHTML(fieldByName('contact_name'),v.contact_name)+fieldHTML(fieldByName('email'),v.email)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('phone'),v.phone)+fieldHTML(fieldByName('category'),v.category)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('status'),v.status||'active')+fieldHTML(fieldByName('contract_end'),v.contract_end)+'</div>';
h+=fieldHTML(fieldByName('annual_spend'),v.annual_spend);
h+=fieldHTML(fieldByName('notes'),v.notes);

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,v[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit){
h+='<button class="btn btn-del" onclick="delVendor()">Delete</button>';
}
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button>';
h+='</div>';
return h;
}

function openForm(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var v=null;
for(var i=0;i<vendors.length;i++){if(vendors[i].id===id){v=vendors[i];break}}
if(!v)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(v);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Company name is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='money')val=parseMoney(el.value);
else if(f.type==='number'||f.type==='integer')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/vendors/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/vendors',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Add failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

async function delVendor(){
if(!editId)return;
if(!confirm('Delete this vendor?'))return;
await fetch(A+'/vendors/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.categories)){
var catField=fieldByName('category');
if(catField)catField.options=cfg.categories;
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
