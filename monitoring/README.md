# Monitoring

Phase 2 introduces a complete local observability baseline.

## Components

- `prometheus/prometheus.yml` — scrape jobs and evaluation setup
- `prometheus/alerts.yml` — starter alert rules
- `otel-collector/config.yaml` — OTLP ingest + Jaeger trace export + Prometheus metric export
- `grafana/provisioning/datasources/datasource.yaml` — Prometheus + Jaeger datasources
- `grafana/provisioning/dashboards/dashboard.yaml` — dashboard auto-loading

## Dashboard files

- `dashboards/grafana/upi-platform-overview.json`
- `dashboards/grafana/kafka-overview.json`

## Access

- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3001` (`admin/admin`)
- Jaeger: `http://localhost:16686`
