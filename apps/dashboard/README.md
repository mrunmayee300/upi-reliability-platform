# UPI Operations Dashboard

Enterprise dark-theme dashboard for the UPI Transaction Intelligence Platform.

## Features

- Real-time metrics via WebSocket
- TPS / latency charts (Recharts)
- Bank health & congestion cards
- Live transaction stream
- Retry & DLQ panel

## Development

```powershell
cd apps/dashboard
copy .env.local.example .env.local
npm install
npm run dev
```

Or from repo root:

```powershell
.\scripts\dev\run-dashboard.ps1
```

**Requires** Analytics service on port **8089** and active traffic from the tx-generator.

## Stack

Next.js 14 · TypeScript · TailwindCSS · Recharts
