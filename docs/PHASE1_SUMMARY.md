# Phase 1 Summary — Architecture & Contracts

## What was delivered

### 1. Monorepo structure

Production-style layout with `apps/`, `services/` (11 + dashboard), `shared/contracts`, `shared/go`, `infra/`, `k8s/`, `helm/`, `terraform/`, `monitoring/`, `load-tests/`, `scripts/`, `docs/`.

### 2. Architecture documentation

- **C4 context** and container diagrams
- **Sequence diagrams:** happy path, retry/DLQ, deployment, routing
- **State machine** for transaction lifecycle
- **Scalability, security, resilience** guides

### 3. Kafka architecture

- 10 topics with partitions, retention, consumer groups
- Standard headers and DLQ pattern
- Partition key strategy per domain entity

### 4. Service contracts

- **AsyncAPI 3.0** for all Kafka channels
- **OpenAPI 3.1** for Gateway, Ingestion, Generator, Bank, Analytics, AI
- **JSON Schema** for envelope, transaction, failure, bank health, congestion, fraud

### 5. Shared Go package

- `models` — transaction, envelope, failure types
- `kafka` — topic and consumer group constants
- `otel` — correlation ID context helpers

### 6. ADRs

Five accepted decisions: Kafka event backbone, Go services, idempotency store, failure pipeline split, OpenTelemetry.

### 7. Database schema (draft)

PostgreSQL tables: transactions, idempotency, retries, DLQ, audit, metrics rollups, fraud alerts.

## Architectural decisions & tradeoffs

| Decision | Why | Tradeoff |
|----------|-----|----------|
| Kafka as SOT for lifecycle | Replay, decouple, fintech-aligned | Ops complexity, eventual UI |
| Go for hot path | Throughput, single binary | Two runtimes with Python ML |
| Detector vs Orchestrator | SRP, independent scale | Consumer coordination care |
| Application-level retry delay | Precise backoff | Not native Kafka delayed topics |
| Redis + PG idempotency | Speed + audit | Redis failure policy required |

## Phase 2 preview

1. `docker-compose.yml` — Kafka (KRaft), Postgres, Redis, Prometheus, Grafana, Jaeger, OTel Collector
2. `scripts/kafka/create-topics.sh`
3. Apply `001_initial_schema.sql`
4. Prometheus scrape configs + starter Grafana dashboards

## Phase 3 preview

Implement services in order: Ingestion → Bank → Failure/Retry → Analytics → Gateway/Generator → Routing/Fraud → Notification → AI.
