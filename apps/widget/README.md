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

Виджет использует нейтральную тёмную тему по умолчанию. Встраивающий фронтенд
может изменить её без форка файла, добавив стили после скрипта:

```css
#hd-root {
  --hd-primary: #fb4b75;
  --hd-primary-foreground: #ffffff;
  --hd-panel: #ffffff;
  --hd-panel-muted: #ffffff;
  --hd-panel-header: #ffffff;
  --hd-panel-bubble: #f4f4f4;
  --hd-panel-border: rgba(35, 35, 35, 0.08);
  --hd-input: #f4f4f4;
  --hd-input-border: transparent;
  --hd-text: #232323;
  --hd-text-muted: #5e5e5e;
  --hd-text-action: #5e5e5e;
  --hd-error: #d94664;
}
```

Local demo: start the stack + web, open `/widget-demo.html`.
