"use strict";var HelpDeskWidget=(()=>{var H=Object.defineProperty;var P=Object.getOwnPropertyDescriptor;var A=Object.getOwnPropertyNames;var R=Object.prototype.hasOwnProperty;var N=(o,s,i,r)=>{if(s&&typeof s=="object"||typeof s=="function")for(let t of A(s))!R.call(o,t)&&t!==i&&H(o,t,{get:()=>s[t],enumerable:!(r=P(s,t))||r.enumerable});return o};var j=o=>N(H({},"__esModule",{value:!0}),o);var J={},_="hd_visitor_id",T="hd_conversation_id";function D(o){let s=o.getAttribute("data-api")||"";if(s)return s.replace(/\/$/,"");let i=o.getAttribute("src")||"";try{let r=new URL(i,window.location.href);return(o.getAttribute("data-gateway")||"http://localhost:8080").replace(/\/$/,"")}catch{return"http://localhost:8080"}}function G(o){return o.getAttribute("data-site-key")||"demo-site"}function U(){return`
#hd-root{all:initial;font-family:ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,sans-serif}
#hd-root *{box-sizing:border-box}
#hd-fab{position:fixed;right:20px;bottom:20px;z-index:2147483000;width:56px;height:56px;border-radius:50%;border:none;cursor:pointer;background:#2563eb;color:#fff;box-shadow:0 8px 24px rgba(37,99,235,.35);font-size:22px;line-height:1}
#hd-fab:hover{filter:brightness(1.05)}
#hd-panel{position:fixed;right:20px;bottom:88px;z-index:2147483000;width:min(380px,calc(100vw - 24px));height:min(520px,calc(100vh - 120px));display:none;flex-direction:column;background:#0f1419;color:#e8eef5;border:1px solid #243041;border-radius:16px;overflow:hidden;box-shadow:0 20px 50px rgba(0,0,0,.45)}
#hd-panel.open{display:flex}
#hd-head{padding:14px 16px;background:#162033;border-bottom:1px solid #243041;display:flex;align-items:center;justify-content:space-between;gap:8px}
#hd-head h2{margin:0;font-size:14px;font-weight:600}
#hd-head button{background:transparent;border:none;color:#9fb0c3;cursor:pointer;font-size:18px}
#hd-msgs{flex:1;overflow:auto;padding:14px;display:flex;flex-direction:column;gap:10px}
.hd-bubble{max-width:85%;padding:10px 12px;border-radius:12px;font-size:13px;line-height:1.45;white-space:pre-wrap;word-break:break-word}
.hd-bubble.visitor{align-self:flex-end;background:#2563eb;color:#fff;border-bottom-right-radius:4px}
.hd-bubble.bot,.hd-bubble.agent,.hd-bubble.system{align-self:flex-start;background:#1c2736;color:#e8eef5;border-bottom-left-radius:4px}
.hd-bubble.system{opacity:.85;font-size:12px}
.hd-meta{font-size:10px;opacity:.65;margin-bottom:4px;text-transform:capitalize}
#hd-foot{border-top:1px solid #243041;padding:10px;display:flex;flex-direction:column;gap:8px;background:#121820}
#hd-row{display:flex;gap:8px}
#hd-input{flex:1;border:1px solid #2a3a4f;background:#0c1118;color:#e8eef5;border-radius:10px;padding:10px 12px;font-size:13px;outline:none}
#hd-input:focus{border-color:#3b82f6}
#hd-send,#hd-human{border:none;border-radius:10px;padding:10px 12px;font-size:12px;font-weight:600;cursor:pointer}
#hd-send{background:#2563eb;color:#fff}
#hd-human{background:#243041;color:#c9d6e5}
#hd-human:disabled,#hd-send:disabled{opacity:.5;cursor:not-allowed}
#hd-err{color:#f87171;font-size:11px;min-height:14px}
`}async function f(o,s,i={}){let r=new Headers(i.headers||{});r.set("Content-Type","application/json"),i.siteKey&&r.set("X-Site-Key",i.siteKey),i.visitorId&&r.set("X-Visitor-Id",i.visitorId);let t=await fetch(`${o}${s}`,{...i,headers:r}),h=await t.json().catch(()=>({}));if(!t.ok)throw new Error(h.error||t.statusText);return h}function C(o){let s=D(o),i=G(o),r=document.createElement("style");r.textContent=U(),document.head.appendChild(r);let t=document.createElement("div");t.id="hd-root",t.innerHTML=`
    <button id="hd-fab" type="button" aria-label="Open support chat">\u{1F4AC}</button>
    <div id="hd-panel" role="dialog" aria-label="Support chat">
      <div id="hd-head">
        <h2>Support</h2>
        <button type="button" id="hd-close" aria-label="Close">\xD7</button>
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
  `,document.body.appendChild(t);let h=t.querySelector("#hd-fab"),y=t.querySelector("#hd-panel"),b=t.querySelector("#hd-msgs"),x=t.querySelector("#hd-input"),v=t.querySelector("#hd-send"),w=t.querySelector("#hd-human"),z=t.querySelector("#hd-err"),q=t.querySelector("#hd-close"),a=localStorage.getItem(_)||"",d=localStorage.getItem(T)||"",c=[],S="",m,p=!1;function u(e){z.textContent=e}function E(){b.innerHTML="";for(let e of c){let n=document.createElement("div");n.className=`hd-bubble ${e.role}`;let l=document.createElement("div");l.className="hd-meta",l.textContent=e.role;let M=document.createElement("div");M.textContent=e.body,n.appendChild(l),n.appendChild(M),b.appendChild(n)}b.scrollTop=b.scrollHeight,c.length&&(S=c[c.length-1].id)}function g(e){let n=new Set(c.map(l=>l.id));for(let l of e)n.has(l.id)||c.push(l);E()}async function O(){if(d&&a)try{c=(await f(s,`/widget/conversations/${d}/messages`,{method:"GET",visitorId:a,siteKey:i})).messages||[],E();return}catch{d="",localStorage.removeItem(T)}let e=await f(s,"/widget/session",{method:"POST",siteKey:i,visitorId:a||void 0,body:JSON.stringify({site_key:i,visitor_id:a||void 0})});a=e.visitor_id,d=e.conversation.id,c=e.messages||[],localStorage.setItem(_,a),localStorage.setItem(T,d),E()}async function $(){if(!(!d||!a||document.hidden))try{let e=S?`?after=${encodeURIComponent(S)}`:"",n=await f(s,`/widget/conversations/${d}/messages${e}`,{method:"GET",visitorId:a,siteKey:i});n.messages?.length&&g(n.messages)}catch{}}function K(){L(),m=window.setInterval($,1500)}function L(){m&&window.clearInterval(m),m=void 0}async function B(){y.classList.add("open"),u("");try{await O(),K()}catch(e){u(e instanceof Error?e.message:"\u041D\u0435 \u0443\u0434\u0430\u043B\u043E\u0441\u044C \u043E\u0442\u043A\u0440\u044B\u0442\u044C \u0447\u0430\u0442")}}function k(){y.classList.remove("open"),L()}h.addEventListener("click",()=>{y.classList.contains("open")?k():B()}),q.addEventListener("click",k);async function I(){let e=x.value.trim();if(!(!e||p||!d)){p=!0,v.disabled=!0,u("");try{let n=await f(s,`/widget/conversations/${d}/messages`,{method:"POST",visitorId:a,siteKey:i,body:JSON.stringify({body:e})});x.value="",n.message&&g([n.message]),n.bot_reply&&g([n.bot_reply])}catch(n){u(n instanceof Error?n.message:"\u041E\u0448\u0438\u0431\u043A\u0430 \u043E\u0442\u043F\u0440\u0430\u0432\u043A\u0438")}finally{p=!1,v.disabled=!1}}}v.addEventListener("click",()=>void I()),x.addEventListener("keydown",e=>{e.key==="Enter"&&I()}),w.addEventListener("click",async()=>{if(!(!d||p)){p=!0,w.disabled=!0,u("");try{let e=await f(s,`/widget/conversations/${d}/handoff`,{method:"POST",visitorId:a,siteKey:i,body:"{}"});e.system_message&&g([e.system_message])}catch(e){u(e instanceof Error?e.message:"Handoff failed")}finally{p=!1,w.disabled=!1}}})}function V(){let o=document.currentScript instanceof HTMLScriptElement?document.currentScript:document.querySelector("script[data-site-key]");o&&(document.readyState==="loading"?document.addEventListener("DOMContentLoaded",()=>C(o)):C(o))}V();return j(J);})();
