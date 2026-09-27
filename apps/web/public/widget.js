"use strict";var HelpDeskWidget=(()=>{var H=Object.defineProperty;var P=Object.getOwnPropertyDescriptor;var A=Object.getOwnPropertyNames;var R=Object.prototype.hasOwnProperty;var N=(n,r,d,i)=>{if(r&&typeof r=="object"||typeof r=="function")for(let t of A(r))!R.call(n,t)&&t!==d&&H(n,t,{get:()=>r[t],enumerable:!(i=P(r,t))||i.enumerable});return n};var j=n=>N(H({},"__esModule",{value:!0}),n);var J={},_="hd_visitor_id",T="hd_conversation_id";function D(n){let r=n.getAttribute("data-api")||"";if(r)return r.replace(/\/$/,"");let d=n.getAttribute("src")||"";try{let i=new URL(d,window.location.href);return(n.getAttribute("data-gateway")||"http://localhost:8080").replace(/\/$/,"")}catch{return"http://localhost:8080"}}function G(n){return n.getAttribute("data-site-key")||"demo-site"}function U(){return`
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
#hd-foot{border-top:1px solid var(--hd-panel-border);padding:10px;display:flex;flex-direction:column;gap:8px;background:var(--hd-panel-muted)}
#hd-row{display:flex;gap:8px}
#hd-input{flex:1;border:1px solid var(--hd-input-border);background:var(--hd-input);color:var(--hd-text);border-radius:10px;padding:10px 12px;font-size:13px;outline:none}
#hd-input:focus{border-color:var(--hd-primary)}
#hd-send,#hd-human{border:none;border-radius:10px;padding:10px 12px;font-size:12px;font-weight:600;cursor:pointer}
#hd-send{background:var(--hd-primary);color:var(--hd-primary-foreground)}
#hd-human{background:var(--hd-panel-border);color:var(--hd-text-action)}
#hd-human:disabled,#hd-send:disabled{opacity:.5;cursor:not-allowed}
#hd-err{color:var(--hd-error);font-size:11px;min-height:14px}
`}async function u(n,r,d={}){let i=new Headers(d.headers||{});i.set("Content-Type","application/json"),d.siteKey&&i.set("X-Site-Key",d.siteKey),d.visitorId&&i.set("X-Visitor-Id",d.visitorId);let t=await fetch(`${n}${r}`,{...d,headers:i}),f=await t.json().catch(()=>({}));if(!t.ok)throw new Error(f.error||t.statusText);return f}function z(n){let r=D(n),d=G(n),i=document.createElement("style");i.textContent=U(),document.head.appendChild(i);let t=document.createElement("div");t.id="hd-root",t.innerHTML=`
    <button id="hd-fab" type="button" aria-label="\u041E\u0442\u043A\u0440\u044B\u0442\u044C \u0447\u0430\u0442 \u043F\u043E\u0434\u0434\u0435\u0440\u0436\u043A\u0438">\u{1F4AC}</button>
    <div id="hd-panel" role="dialog" aria-label="\u0427\u0430\u0442 \u043F\u043E\u0434\u0434\u0435\u0440\u0436\u043A\u0438">
      <div id="hd-head">
        <h2>\u041F\u043E\u0434\u0434\u0435\u0440\u0436\u043A\u0430</h2>
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
  `,document.body.appendChild(t);let f=t.querySelector("#hd-fab"),y=t.querySelector("#hd-panel"),b=t.querySelector("#hd-msgs"),x=t.querySelector("#hd-input"),v=t.querySelector("#hd-send"),w=t.querySelector("#hd-human"),C=t.querySelector("#hd-err"),q=t.querySelector("#hd-close"),s=localStorage.getItem(_)||"",a=localStorage.getItem(T)||"",l=[],E="",m,p=!1;function h(e){C.textContent=e}function S(){b.innerHTML="";for(let e of l){let o=document.createElement("div");o.className=`hd-bubble ${e.role}`;let c=document.createElement("div");c.className="hd-meta",c.textContent=e.role;let M=document.createElement("div");M.textContent=e.body,o.appendChild(c),o.appendChild(M),b.appendChild(o)}b.scrollTop=b.scrollHeight,l.length&&(E=l[l.length-1].id)}function g(e){let o=new Set(l.map(c=>c.id));for(let c of e)o.has(c.id)||l.push(c);S()}async function $(){if(a&&s)try{l=(await u(r,`/widget/conversations/${a}/messages`,{method:"GET",visitorId:s,siteKey:d})).messages||[],S();return}catch{a="",localStorage.removeItem(T)}let e=await u(r,"/widget/session",{method:"POST",siteKey:d,visitorId:s||void 0,body:JSON.stringify({site_key:d,visitor_id:s||void 0})});s=e.visitor_id,a=e.conversation.id,l=e.messages||[],localStorage.setItem(_,s),localStorage.setItem(T,a),S()}async function K(){if(!(!a||!s||document.hidden))try{let e=E?`?after=${encodeURIComponent(E)}`:"",o=await u(r,`/widget/conversations/${a}/messages${e}`,{method:"GET",visitorId:s,siteKey:d});o.messages?.length&&g(o.messages)}catch{}}function O(){L(),m=window.setInterval(K,1500)}function L(){m&&window.clearInterval(m),m=void 0}async function B(){y.classList.add("open"),h("");try{await $(),O()}catch(e){h(e instanceof Error?e.message:"\u041D\u0435 \u0443\u0434\u0430\u043B\u043E\u0441\u044C \u043E\u0442\u043A\u0440\u044B\u0442\u044C \u0447\u0430\u0442")}}function k(){y.classList.remove("open"),L()}f.addEventListener("click",()=>{y.classList.contains("open")?k():B()}),q.addEventListener("click",k);async function I(){let e=x.value.trim();if(!(!e||p||!a)){p=!0,v.disabled=!0,h("");try{let o=await u(r,`/widget/conversations/${a}/messages`,{method:"POST",visitorId:s,siteKey:d,body:JSON.stringify({body:e})});x.value="",o.message&&g([o.message]),o.bot_reply&&g([o.bot_reply])}catch(o){h(o instanceof Error?o.message:"\u041E\u0448\u0438\u0431\u043A\u0430 \u043E\u0442\u043F\u0440\u0430\u0432\u043A\u0438")}finally{p=!1,v.disabled=!1}}}v.addEventListener("click",()=>void I()),x.addEventListener("keydown",e=>{e.key==="Enter"&&I()}),w.addEventListener("click",async()=>{if(!(!a||p)){p=!0,w.disabled=!0,h("");try{let e=await u(r,`/widget/conversations/${a}/handoff`,{method:"POST",visitorId:s,siteKey:d,body:"{}"});e.system_message&&g([e.system_message])}catch(e){h(e instanceof Error?e.message:"Handoff failed")}finally{p=!1,w.disabled=!1}}})}function V(){let n=document.currentScript instanceof HTMLScriptElement?document.currentScript:document.querySelector("script[data-site-key]");n&&(document.readyState==="loading"?document.addEventListener("DOMContentLoaded",()=>z(n)):z(n))}V();return j(J);})();
