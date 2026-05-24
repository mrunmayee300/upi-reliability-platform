# Documentation Images

SVG diagrams for README and architecture docs. Render natively on GitHub.

| File | Description | Used in |
|------|-------------|---------|
| [banner.svg](./banner.svg) | Project header banner | README |
| [architecture-overview.svg](./architecture-overview.svg) | Full platform architecture | README, architecture docs |
| [transaction-flow.svg](./transaction-flow.svg) | Happy path + failure/recovery | README, data-flow |
| [deployment-topology.svg](./deployment-topology.svg) | Dev / Compose / K8s options | README, deployment runbook |
| [dashboard-preview.svg](./dashboard-preview.svg) | Operations dashboard mockup | README |
| [observability-stack.svg](./observability-stack.svg) | Metrics, traces, dashboards | README, architecture |

### Adding a real dashboard screenshot

After running the dashboard locally:

```powershell
# Start stack, open http://localhost:3000, capture screenshot
# Save as docs/images/dashboard-screenshot.png
```

Then reference in README:

```markdown
![Dashboard screenshot](./docs/images/dashboard-screenshot.png)
```
