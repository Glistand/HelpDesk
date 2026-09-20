# Embeddable support chat widget

Build copies the IIFE bundle to `apps/web/public/widget.js`.

```bash
npm install
npm run build
```

Embed:

```html
<script
  src="https://your-host/widget.js"
  data-site-key="demo-site"
  data-gateway="https://api.example.com"
></script>
```

Local demo: start the stack + web, open `/widget-demo.html`.
