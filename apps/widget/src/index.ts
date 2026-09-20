type Msg = {
  id: string;
  role: string;
  body: string;
  created_at: string;
};

type Session = {
  visitor_id: string;
  conversation: { id: string; status: string };
  messages: Msg[];
};

const STORAGE_VISITOR = "hd_visitor_id";
const STORAGE_CONV = "hd_conversation_id";

function apiBase(script: HTMLScriptElement): string {
  const fromAttr = script.getAttribute("data-api") || "";
  if (fromAttr) return fromAttr.replace(/\/$/, "");
  const src = script.getAttribute("src") || "";
  try {
    const u = new URL(src, window.location.href);
    // widget.js served from web; gateway is typically :8080
    return (script.getAttribute("data-gateway") || "http://localhost:8080").replace(/\/$/, "");
  } catch {
    return "http://localhost:8080";
  }
}

function siteKey(script: HTMLScriptElement): string {
  return script.getAttribute("data-site-key") || "demo-site";
}

function css(): string {
  return `
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
`;
}

async function jsonFetch(
  base: string,
  path: string,
  init: RequestInit & { visitorId?: string; siteKey?: string } = {},
): Promise<any> {
  const headers = new Headers(init.headers || {});
  headers.set("Content-Type", "application/json");
  if (init.siteKey) headers.set("X-Site-Key", init.siteKey);
  if (init.visitorId) headers.set("X-Visitor-Id", init.visitorId);
  const res = await fetch(`${base}${path}`, { ...init, headers });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || res.statusText);
  return data;
}

function mount(script: HTMLScriptElement) {
  const base = apiBase(script);
  const key = siteKey(script);

  const style = document.createElement("style");
  style.textContent = css();
  document.head.appendChild(style);

  const root = document.createElement("div");
  root.id = "hd-root";
  root.innerHTML = `
    <button id="hd-fab" type="button" aria-label="Open support chat">💬</button>
    <div id="hd-panel" role="dialog" aria-label="Support chat">
      <div id="hd-head">
        <h2>Support</h2>
        <button type="button" id="hd-close" aria-label="Close">×</button>
      </div>
      <div id="hd-msgs"></div>
      <div id="hd-foot">
        <div id="hd-err"></div>
        <div id="hd-row">
          <input id="hd-input" type="text" placeholder="Напишите сообщение…" autocomplete="off" />
          <button id="hd-send" type="button">Отправить</button>
        </div>
        <button id="hd-human" type="button">Нужен человек</button>
      </div>
    </div>
  `;
  document.body.appendChild(root);

  const fab = root.querySelector("#hd-fab") as HTMLButtonElement;
  const panel = root.querySelector("#hd-panel") as HTMLDivElement;
  const msgsEl = root.querySelector("#hd-msgs") as HTMLDivElement;
  const input = root.querySelector("#hd-input") as HTMLInputElement;
  const sendBtn = root.querySelector("#hd-send") as HTMLButtonElement;
  const humanBtn = root.querySelector("#hd-human") as HTMLButtonElement;
  const errEl = root.querySelector("#hd-err") as HTMLDivElement;
  const closeBtn = root.querySelector("#hd-close") as HTMLButtonElement;

  let visitorId = localStorage.getItem(STORAGE_VISITOR) || "";
  let conversationId = localStorage.getItem(STORAGE_CONV) || "";
  let messages: Msg[] = [];
  let lastId = "";
  let pollTimer: number | undefined;
  let busy = false;

  function setErr(msg: string) {
    errEl.textContent = msg;
  }

  function render() {
    msgsEl.innerHTML = "";
    for (const m of messages) {
      const wrap = document.createElement("div");
      wrap.className = `hd-bubble ${m.role}`;
      const meta = document.createElement("div");
      meta.className = "hd-meta";
      meta.textContent = m.role;
      const body = document.createElement("div");
      body.textContent = m.body;
      wrap.appendChild(meta);
      wrap.appendChild(body);
      msgsEl.appendChild(wrap);
    }
    msgsEl.scrollTop = msgsEl.scrollHeight;
    if (messages.length) lastId = messages[messages.length - 1].id;
  }

  function mergeMessages(incoming: Msg[]) {
    const seen = new Set(messages.map((m) => m.id));
    for (const m of incoming) {
      if (!seen.has(m.id)) messages.push(m);
    }
    render();
  }

  async function ensureSession() {
    if (conversationId && visitorId) {
      try {
        const data = await jsonFetch(base, `/widget/conversations/${conversationId}/messages`, {
          method: "GET",
          visitorId,
          siteKey: key,
        });
        messages = data.messages || [];
        render();
        return;
      } catch {
        conversationId = "";
        localStorage.removeItem(STORAGE_CONV);
      }
    }
    const data: Session = await jsonFetch(base, "/widget/session", {
      method: "POST",
      siteKey: key,
      visitorId: visitorId || undefined,
      body: JSON.stringify({
        site_key: key,
        visitor_id: visitorId || undefined,
      }),
    });
    visitorId = data.visitor_id;
    conversationId = data.conversation.id;
    messages = data.messages || [];
    localStorage.setItem(STORAGE_VISITOR, visitorId);
    localStorage.setItem(STORAGE_CONV, conversationId);
    render();
  }

  async function poll() {
    if (!conversationId || !visitorId || document.hidden) return;
    try {
      const q = lastId ? `?after=${encodeURIComponent(lastId)}` : "";
      const data = await jsonFetch(base, `/widget/conversations/${conversationId}/messages${q}`, {
        method: "GET",
        visitorId,
        siteKey: key,
      });
      if (data.messages?.length) mergeMessages(data.messages);
    } catch {
      /* ignore transient poll errors */
    }
  }

  function startPoll() {
    stopPoll();
    pollTimer = window.setInterval(poll, 1500);
  }

  function stopPoll() {
    if (pollTimer) window.clearInterval(pollTimer);
    pollTimer = undefined;
  }

  async function openPanel() {
    panel.classList.add("open");
    setErr("");
    try {
      await ensureSession();
      startPoll();
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Не удалось открыть чат");
    }
  }

  function closePanel() {
    panel.classList.remove("open");
    stopPoll();
  }

  fab.addEventListener("click", () => {
    if (panel.classList.contains("open")) closePanel();
    else void openPanel();
  });
  closeBtn.addEventListener("click", closePanel);

  async function send() {
    const body = input.value.trim();
    if (!body || busy || !conversationId) return;
    busy = true;
    sendBtn.disabled = true;
    setErr("");
    try {
      const data = await jsonFetch(base, `/widget/conversations/${conversationId}/messages`, {
        method: "POST",
        visitorId,
        siteKey: key,
        body: JSON.stringify({ body }),
      });
      input.value = "";
      if (data.message) mergeMessages([data.message]);
      if (data.bot_reply) mergeMessages([data.bot_reply]);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Ошибка отправки");
    } finally {
      busy = false;
      sendBtn.disabled = false;
    }
  }

  sendBtn.addEventListener("click", () => void send());
  input.addEventListener("keydown", (e) => {
    if (e.key === "Enter") void send();
  });

  humanBtn.addEventListener("click", async () => {
    if (!conversationId || busy) return;
    busy = true;
    humanBtn.disabled = true;
    setErr("");
    try {
      const data = await jsonFetch(base, `/widget/conversations/${conversationId}/handoff`, {
        method: "POST",
        visitorId,
        siteKey: key,
        body: "{}",
      });
      if (data.system_message) mergeMessages([data.system_message]);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Handoff failed");
    } finally {
      busy = false;
      humanBtn.disabled = false;
    }
  });
}

function boot() {
  const script =
    document.currentScript instanceof HTMLScriptElement
      ? document.currentScript
      : (document.querySelector("script[data-site-key]") as HTMLScriptElement | null);
  if (!script) return;
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", () => mount(script));
  } else {
    mount(script);
  }
}

boot();

export {};
