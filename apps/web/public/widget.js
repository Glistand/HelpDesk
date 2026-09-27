"use strict";var HelpDeskWidget=(()=>{var $=Object.defineProperty;var R=Object.getOwnPropertyDescriptor;var B=Object.getOwnPropertyNames;var U=Object.prototype.hasOwnProperty;var j=(t,o,r,d)=>{if(o&&typeof o=="object"||typeof o=="function")for(let e of B(o))!U.call(t,e)&&e!==r&&$(t,e,{get:()=>o[e],enumerable:!(d=R(o,e))||d.enumerable});return t};var D=t=>j($({},"__esModule",{value:!0}),t);var W={},q="hd_visitor_id",H="hd_conversation_id";function G(t){let o=t.getAttribute("data-api")||"";if(o)return o.replace(/\/$/,"");let r=t.getAttribute("src")||"";try{let d=new URL(r,window.location.href);return(t.getAttribute("data-gateway")||"http://localhost:8080").replace(/\/$/,"")}catch{return"http://localhost:8080"}}function V(t){return t.getAttribute("data-site-key")||"demo-site"}function F(t){return t.getAttribute("data-title")||"\u041F\u043E\u0434\u0434\u0435\u0440\u0436\u043A\u0430"}function J(){return`
#hd-root{all:initial;font-family:ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,sans-serif;--hd-primary:#2563eb;--hd-primary-foreground:#fff;--hd-panel:#0f1419;--hd-panel-muted:#121820;--hd-panel-header:#162033;--hd-panel-bubble:#1c2736;--hd-panel-border:#243041;--hd-input:#0c1118;--hd-input-border:#2a3a4f;--hd-text:#e8eef5;--hd-text-muted:#9fb0c3;--hd-text-action:#c9d6e5;--hd-error:#f87171}
#hd-root *{box-sizing:border-box}
#hd-fab{position:fixed;right:20px;bottom:20px;z-index:2147483000;width:56px;height:56px;border-radius:50%;border:none;cursor:pointer;background:var(--hd-primary);color:var(--hd-primary-foreground);box-shadow:0 8px 24px color-mix(in srgb,var(--hd-primary) 35%,transparent);font-size:22px;line-height:1}
#hd-fab:hover{filter:brightness(1.05)}
#hd-fab:focus-visible,#hd-send:focus-visible,#hd-human:focus-visible,#hd-close:focus-visible,#hd-input:focus-visible{outline:2px solid var(--hd-primary);outline-offset:2px}
#hd-panel{position:fixed;right:20px;bottom:88px;z-index:2147483000;width:min(380px,calc(100vw - 24px));height:min(520px,calc(100vh - 120px));display:none;flex-direction:column;background:var(--hd-panel);color:var(--hd-text);border:1px solid var(--hd-panel-border);border-radius:16px;overflow:hidden;box-shadow:0 20px 50px rgba(0,0,0,.45)}
#hd-panel.open{display:flex}
#hd-head{padding:14px 16px;background:var(--hd-panel-header);border-bottom:1px solid var(--hd-panel-border);display:flex;align-items:center;justify-content:space-between;gap:8px}
#hd-head h2{margin:0;font-size:14px;font-weight:600}
#hd-head button{background:transparent;border:none;color:var(--hd-text-muted);cursor:pointer;font-size:18px}
#hd-msgs{flex:1;overflow:auto;padding:14px;display:flex;flex-direction:column;gap:10px}
.hd-bubble{max-width:85%;padding:10px 12px;border-radius:12px;font-size:13px;line-height:1.45;white-space:pre-wrap;word-break:break-word}
.hd-bubble.visitor{align-self:flex-end;background:var(--hd-primary);color:var(--hd-primary-foreground);border-bottom-right-radius:4px}
.hd-bubble.bot,.hd-bubble.agent,.hd-bubble.system{align-self:flex-start;background:var(--hd-panel-bubble);color:var(--hd-text);border-bottom-left-radius:4px}
.hd-bubble.system{opacity:.85;font-size:12px}
.hd-meta{font-size:10px;opacity:.65;margin-bottom:4px;text-transform:capitalize}
.hd-markdown p{margin:0 0 8px}
.hd-markdown p:last-child{margin-bottom:0}
.hd-markdown h1,.hd-markdown h2,.hd-markdown h3{margin:0 0 8px;font-size:1em;line-height:1.35}
.hd-markdown ul,.hd-markdown ol{margin:0 0 8px;padding-left:18px}
.hd-markdown li{margin:3px 0}
.hd-markdown strong{font-weight:700}
.hd-markdown em{font-style:italic}
.hd-markdown code{padding:1px 4px;border-radius:4px;background:rgba(127,127,127,.16);font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:.92em}
.hd-markdown a{color:inherit;text-decoration:underline;text-underline-offset:2px}
#hd-foot{border-top:1px solid var(--hd-panel-border);padding:10px;display:flex;flex-direction:column;gap:8px;background:var(--hd-panel-muted)}
#hd-row{display:flex;gap:8px}
#hd-input{flex:1;border:1px solid var(--hd-input-border);background:var(--hd-input);color:var(--hd-text);border-radius:10px;padding:10px 12px;font-size:13px;outline:none}
#hd-input:focus{border-color:var(--hd-primary)}
#hd-send,#hd-human{border:none;border-radius:10px;padding:10px 12px;font-size:12px;font-weight:600;cursor:pointer}
#hd-send{background:var(--hd-primary);color:var(--hd-primary-foreground)}
#hd-human{background:var(--hd-panel-border);color:var(--hd-text-action)}
#hd-human:disabled,#hd-send:disabled{opacity:.5;cursor:not-allowed}
#hd-err{color:var(--hd-error);font-size:11px;min-height:14px}
`}async function v(t,o,r={}){let d=new Headers(r.headers||{});d.set("Content-Type","application/json"),r.siteKey&&d.set("X-Site-Key",r.siteKey),r.visitorId&&d.set("X-Visitor-Id",r.visitorId);let e=await fetch(`${t}${o}`,{...r,headers:d}),c=await e.json().catch(()=>({}));if(!e.ok)throw new Error(c.error||e.statusText);return c}function I(t,o){let r=/(\*\*([^*]+)\*\*|`([^`]+)`|\[([^\]]+)\]\(([^)\s]+)\)|\*([^*]+)\*)/g,d=0;for(let e of o.matchAll(r)){let c=e.index??0;if(c>d&&t.append(document.createTextNode(o.slice(d,c))),e[2]){let i=document.createElement("strong");i.textContent=e[2],t.append(i)}else if(e[3]){let i=document.createElement("code");i.textContent=e[3],t.append(i)}else if(e[4]&&e[5])try{let i=new URL(e[5],window.location.href);if(i.protocol==="http:"||i.protocol==="https:"){let s=document.createElement("a");s.href=i.toString(),s.target="_blank",s.rel="noopener noreferrer",s.textContent=e[4],t.append(s)}else t.append(document.createTextNode(e[0]))}catch{t.append(document.createTextNode(e[0]))}else if(e[6]){let i=document.createElement("em");i.textContent=e[6],t.append(i)}d=c+e[0].length}d<o.length&&t.append(document.createTextNode(o.slice(d)))}function X(t,o){let r=o.replace(/\r\n?/g,`
`).split(`
`),d=[],e,c,i=()=>{if(!d.length)return;let l=document.createElement("p");d.forEach((u,m)=>{m&&l.append(document.createElement("br")),I(l,u)}),t.append(l),d=[]},s=()=>{e=void 0,c=void 0};for(let l of r){let u=/^(#{1,3})\s+(.+)$/.exec(l),m=/^\s*[-+*]\s+(.+)$/.exec(l),w=/^\s*\d+[.)]\s+(.+)$/.exec(l);if(!l.trim()){i(),s();continue}if(u){i(),s();let f=document.createElement(`h${u[1].length}`);I(f,u[2]),t.append(f);continue}if(m||w){i();let f=m?"ul":"ol";(!e||c!==f)&&(e=document.createElement(f),c=f,t.append(e));let E=document.createElement("li");I(E,(m??w)[1]),e.append(E);continue}s(),d.push(l)}i()}function K(t){let o=G(t),r=V(t),d=document.createElement("style");d.textContent=J(),document.head.appendChild(d);let e=document.createElement("div");e.id="hd-root",e.innerHTML=`
    <button id="hd-fab" type="button" aria-label="\u041E\u0442\u043A\u0440\u044B\u0442\u044C \u0447\u0430\u0442 \u043F\u043E\u0434\u0434\u0435\u0440\u0436\u043A\u0438">\u{1F4AC}</button>
    <div id="hd-panel" role="dialog" aria-label="\u0427\u0430\u0442 \u043F\u043E\u0434\u0434\u0435\u0440\u0436\u043A\u0438">
      <div id="hd-head">
        <h2 id="hd-title"></h2>
        <button type="button" id="hd-close" aria-label="\u0417\u0430\u043A\u0440\u044B\u0442\u044C">\xD7</button>
      </div>
      <div id="hd-msgs"></div>
      <div id="hd-foot">
        <div id="hd-err"></div>
        <div id="hd-row">
          <input id="hd-input" type="text" placeholder="\u041D\u0430\u043F\u0438\u0448\u0438\u0442\u0435 \u0441\u043E\u043E\u0431\u0449\u0435\u043D\u0438\u0435\u2026" autocomplete="off" />
          <button id="hd-send" type="button">\u041E\u0442\u043F\u0440\u0430\u0432\u0438\u0442\u044C</button>
        </div>
        <button id="hd-human" type="button">\u041D\u0443\u0436\u0435\u043D \u0447\u0435\u043B\u043E\u0432\u0435\u043A</button>
      </div>
    </div>
  `,document.body.appendChild(e);let c=e.querySelector("#hd-fab"),i=e.querySelector("#hd-panel"),s=e.querySelector("#hd-msgs"),l=e.querySelector("#hd-input"),u=e.querySelector("#hd-send"),m=e.querySelector("#hd-human"),w=e.querySelector("#hd-err"),f=e.querySelector("#hd-close"),E=e.querySelector("#hd-title");E.textContent=F(t);let h=localStorage.getItem(q)||"",p=localStorage.getItem(H)||"",b=[],T="",k,x=!1;function y(n){w.textContent=n}function L(){s.innerHTML="";for(let n of b){let a=document.createElement("div");a.className=`hd-bubble ${n.role}`;let g=document.createElement("div");g.className="hd-meta",g.textContent=n.role;let M=document.createElement("div");M.className="hd-markdown",X(M,n.body),a.appendChild(g),a.appendChild(M),s.appendChild(a)}s.scrollTop=s.scrollHeight,b.length&&(T=b[b.length-1].id)}function S(n){let a=new Set(b.map(g=>g.id));for(let g of n)a.has(g.id)||b.push(g);L()}async function O(){if(p&&h)try{b=(await v(o,`/widget/conversations/${p}/messages`,{method:"GET",visitorId:h,siteKey:r})).messages||[],L();return}catch{p="",localStorage.removeItem(H)}let n=await v(o,"/widget/session",{method:"POST",siteKey:r,visitorId:h||void 0,body:JSON.stringify({site_key:r,visitor_id:h||void 0})});h=n.visitor_id,p=n.conversation.id,b=n.messages||[],localStorage.setItem(q,h),localStorage.setItem(H,p),L()}async function N(){if(!(!p||!h||document.hidden))try{let n=T?`?after=${encodeURIComponent(T)}`:"",a=await v(o,`/widget/conversations/${p}/messages${n}`,{method:"GET",visitorId:h,siteKey:r});a.messages?.length&&S(a.messages)}catch{}}function A(){C(),k=window.setInterval(N,1500)}function C(){k&&window.clearInterval(k),k=void 0}async function P(){i.classList.add("open"),y("");try{await O(),A()}catch(n){y(n instanceof Error?n.message:"\u041D\u0435 \u0443\u0434\u0430\u043B\u043E\u0441\u044C \u043E\u0442\u043A\u0440\u044B\u0442\u044C \u0447\u0430\u0442")}}function _(){i.classList.remove("open"),C()}c.addEventListener("click",()=>{i.classList.contains("open")?_():P()}),f.addEventListener("click",_);async function z(){let n=l.value.trim();if(!(!n||x||!p)){x=!0,u.disabled=!0,y("");try{let a=await v(o,`/widget/conversations/${p}/messages`,{method:"POST",visitorId:h,siteKey:r,body:JSON.stringify({body:n})});l.value="",a.message&&S([a.message]),a.bot_reply&&S([a.bot_reply])}catch(a){y(a instanceof Error?a.message:"\u041E\u0448\u0438\u0431\u043A\u0430 \u043E\u0442\u043F\u0440\u0430\u0432\u043A\u0438")}finally{x=!1,u.disabled=!1}}}u.addEventListener("click",()=>void z()),l.addEventListener("keydown",n=>{n.key==="Enter"&&z()}),m.addEventListener("click",async()=>{if(!(!p||x)){x=!0,m.disabled=!0,y("");try{let n=await v(o,`/widget/conversations/${p}/handoff`,{method:"POST",visitorId:h,siteKey:r,body:"{}"});n.system_message&&S([n.system_message])}catch(n){y(n instanceof Error?n.message:"Handoff failed")}finally{x=!1,m.disabled=!1}}})}function Q(){let t=document.currentScript instanceof HTMLScriptElement?document.currentScript:document.querySelector("script[data-site-key]");t&&(document.readyState==="loading"?document.addEventListener("DOMContentLoaded",()=>K(t)):K(t))}Q();return D(W);})();
