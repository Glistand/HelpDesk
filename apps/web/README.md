# Support Desk — agent console

Next.js UI for conversation inbox/thread. Talks to `api-gateway` via cookie `hd_token` and `/api/hd/*` proxy.

```bash
cp .env.example .env.local   # GATEWAY_URL=http://localhost:8080
npm run dev -- --port 3001
```

- Login: `agent@helpdesk.local` / `password`
- Inbox: `/inbox`
- Widget demo (static): `/widget-demo.html`
- Rebuild widget: `cd ../widget && npm run build`
