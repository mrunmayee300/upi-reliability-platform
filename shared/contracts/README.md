# Service Contracts

Canonical API and event contracts for the UPI platform.

## Layout

```
shared/contracts/
├── asyncapi/asyncapi.yaml      # Kafka channels & messages
├── openapi/                    # REST per microservice
│   ├── api-gateway.yaml
│   ├── tx-ingestion.yaml
│   ├── tx-generator.yaml
│   ├── bank-simulator.yaml
│   ├── analytics.yaml
│   └── ai-prediction.yaml
└── schemas/                    # JSON Schema payloads
    ├── event-envelope.json
    ├── transaction.json
    ├── transaction-failure.json
    ├── bank-health.json
    ├── congestion-event.json
    └── fraud-alert.json
```

## Service port map (local)

| Service | Port | Contract |
|---------|------|----------|
| API Gateway | 8080 | `openapi/api-gateway.yaml` |
| Tx Ingestion | 8081 | `openapi/tx-ingestion.yaml` |
| Tx Generator | 8082 | `openapi/tx-generator.yaml` |
| Bank Simulator | 8083 | `openapi/bank-simulator.yaml` |
| Failure Detector | 8084 | Kafka only |
| Retry Orchestrator | 8085 | Kafka + admin REST (Phase 3) |
| Intelligent Routing | 8086 | Kafka only |
| Fraud Detector | 8087 | Kafka only |
| Analytics | 8089 | `openapi/analytics.yaml` |
| Notification | 8090 | Kafka + webhook REST (Phase 3) |
| AI Prediction | 8091 | `openapi/ai-prediction.yaml` |
| Dashboard | 3000 | Consumes Analytics WS |

## Versioning

- `event_version` in envelope: semver minor (`1.0`, `1.1`)
- Breaking changes → new topic suffix or consumer group version (`-v2`)

## Code generation (Phase 3)

Go types generated from JSON Schema via `go-jsonschema` or hand-maintained in `shared/go/models`.
