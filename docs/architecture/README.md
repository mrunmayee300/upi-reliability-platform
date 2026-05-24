# Architecture — UPI Transaction Intelligence & Failure Recovery Platform

## Document index

| Document | Description |
|----------|-------------|
| [system-overview.md](./system-overview.md) | End-to-end flows, boundaries, consistency model |
| [diagrams.md](./diagrams.md) | Architecture images + Mermaid (C4, sequence, deployment) |
| [../images/](../images/README.md) | SVG diagrams (banner, architecture, flows, dashboard) |
| [data-flow.md](./data-flow.md) | Event-driven pipelines and state transitions |
| [scalability.md](./scalability.md) | Partitioning, autoscaling, capacity planning |
| [security.md](./security.md) | Auth simulation, secrets, network policies |
| [resilience.md](./resilience.md) | Retries, circuit breakers, DLQ, chaos |

## Design principles

1. **Event-first** — Kafka is the system of record for in-flight payment lifecycle; PostgreSQL for durable aggregates and audit.
2. **At-least-once with idempotency** — Consumers dedupe via `idempotency_key` + Redis/DB idempotency store.
3. **Failure as a first-class domain** — Failed payments are events, not log lines; recovery is orchestrated, not ad-hoc.
4. **Observable by default** — Correlation IDs propagate gateway → ingestion → bank → retry; OTel spans tie HTTP and Kafka.
5. **Bank as a simulated dependency** — Congestion, latency, and outage are modeled explicitly for routing and prediction.

## Service map

![Platform architecture](../images/architecture-overview.svg)

## Phase roadmap

| Phase | Scope | Status |
|-------|--------|--------|
| 1 | Architecture, structure, contracts | Complete |
| 2 | Infra: Kafka, Postgres, Redis, observability | Complete |
| 3 | Go microservices implementation | Complete |
| 4 | Next.js dashboard | Complete |
| 5 | K8s, Helm, Terraform | Complete |
| 6 | k6 load tests, chaos engineering | Complete |
| 7 | Production README, runbooks, benchmarks | Complete |

## Key tradeoffs

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Primary language | Go (Gin) | Low latency, strong concurrency, fits payment gateways |
| Stream bus | Kafka | Industry standard for payment event logs; replay + ordering |
| Stream processing | Flink (optional path) | Complex windowed analytics; Phase 2+ for congestion windows |
| Sync API | REST + WebSocket | Dashboard and gateway simplicity; gRPC internal optional later |
| Idempotency store | Redis + Postgres | Redis for hot dedupe; Postgres for audit and DLQ correlation |
| ML service | Python (FastAPI + XGBoost) | Ecosystem for time-series forecasting |

See [../adr/](../adr/) for recorded architecture decisions.
