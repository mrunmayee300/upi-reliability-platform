# Documentation Images

| File | Description |
|------|-------------|
| [architecture-overview.png](./architecture-overview.png) | Full platform architecture (PNG) |
| [architecture-overview.svg](./architecture-overview.svg) | Source SVG |
| [transaction-flow.png](./transaction-flow.png) | Happy path + failure/recovery (PNG) |
| [transaction-flow.svg](./transaction-flow.svg) | Source SVG |
| [banner.svg](./banner.svg) | Project banner |
| [dashboard-preview.svg](./dashboard-preview.svg) | Dashboard mockup |
| [deployment-topology.svg](./deployment-topology.svg) | Deploy options |
| [observability-stack.svg](./observability-stack.svg) | Observability pipeline |

### Regenerate PNG from SVG

```powershell
cd C:\Users\Mrunmayee\OneDrive\Desktop\upi
node scripts/dev/svg-to-png.mjs
```

Requires `@resvg/resvg-js` (installed on first run via `npm install --no-save @resvg/resvg-js`).
