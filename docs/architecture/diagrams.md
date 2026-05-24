# Architecture Diagrams

Static diagrams (render on GitHub without Mermaid):

| Diagram | Image |
|---------|-------|
| Platform overview | ![Architecture](../images/architecture-overview.svg) |
| Transaction flow | ![Transaction flow](../images/transaction-flow.svg) |
| Deployment options | ![Deployment](../images/deployment-topology.svg) |
| Observability | ![Observability](../images/observability-stack.svg) |
| Dashboard UI | ![Dashboard](../images/dashboard-preview.svg) |

Full set: [docs/images/README.md](../images/README.md)

---

## Container diagram (C4)

```mermaid
flowchart TB
    subgraph clients [Clients]
        UI[Dashboard Next.js]
        GEN[Tx Generator CLI/API]
    end

    subgraph edge [Edge]
        GW[API Gateway]
    end

    subgraph ingest [Ingestion Path]
        ING[Tx Ingestion]
        BANK[Bank Simulator]
    end

    subgraph process [Processing Path]
        FAIL[Failure Detector]
        RETRY[Retry Orchestrator]
        ROUTE[Intelligent Routing]
        FRAUD[Fraud Detector]
        ANAL[Analytics]
    end

    subgraph intel [Intelligence]
        AI[AI Prediction Python]
        NOTIFY[Notification]
    end

    subgraph data [Data Plane]
        KAFKA[(Kafka)]
        PG[(PostgreSQL)]
        REDIS[(Redis)]
    end

    subgraph obs [Observability]
        PROM[Prometheus]
        GRAF[Grafana]
        JAEG[Jaeger]
    end

    UI --> GW
    GEN --> ING
    GW --> ING
    GW --> ANAL
  UI -->|WebSocket| ANAL

    ING --> KAFKA
    BANK --> KAFKA
    FAIL --> KAFKA
    RETRY --> KAFKA
    ROUTE --> KAFKA
    FRAUD --> KAFKA
    ANAL --> KAFKA

    ING --> PG
    RETRY --> REDIS
    RETRY --> PG
    ROUTE --> REDIS
    ANAL --> PG
    FRAUD --> PG

    AI --> ANAL
    AI --> KAFKA
    NOTIFY --> KAFKA

    GW -.-> JAEG
    ING -.-> PROM
    BANK -.-> PROM
```

## Happy-path sequence

```mermaid
sequenceDiagram
    autonumber
    participant G as Tx Generator
    participant I as Tx Ingestion
    participant K as Kafka
    participant B as Bank Simulator
    participant A as Analytics
    participant D as Dashboard

    G->>I: POST /v1/transactions
    I->>I: Validate + idempotency check
    I->>K: upi-transactions
    I->>K: validated-transactions
    K->>B: consume validated
    B->>B: Simulate bank latency
    B->>K: latency-events
    B->>K: bank-health (periodic)
    alt Success
        B->>K: analytics-events (SUCCESS)
    else Failure
        B->>K: failed-transactions
        K->>FAIL: Failure Detector
        Note over FAIL: Classify failure
        K->>RETRY: Retry Orchestrator
        RETRY->>K: retry-transactions
    end
    K->>A: analytics-events
    A->>D: WebSocket metrics push
```

## Retry & DLQ sequence

```mermaid
sequenceDiagram
    participant K as Kafka
    participant R as Retry Orchestrator
    participant Redis as Redis Idempotency
    participant PG as PostgreSQL
    participant B as Bank Simulator

    K->>R: failed-transactions
    R->>Redis: CHECK idempotency_key
    alt Already processed
        R-->>K: skip (duplicate)
    else New retry
        R->>R: Compute backoff delay
        R->>PG: INSERT retry_attempt
        R->>K: retry-transactions (scheduled)
    end
    Note over R: Max attempts exceeded
    R->>K: dead-letter-events
    R->>K: analytics-events (DLQ)
```

## Deployment view (Kubernetes)

```mermaid
flowchart TB
    subgraph ns_upi [Namespace: upi-platform]
        ING_POD[tx-ingestion x N]
        BANK_POD[bank-simulator x N]
        RETRY_POD[retry-orchestrator x N]
        HPA[HPA]
    end

    subgraph ns_kafka [Namespace: kafka]
        KAFKA_STRIMZI[Kafka Cluster]
    end

    subgraph ns_obs [Namespace: observability]
        PROM_ST[Prometheus]
        GRAF_ST[Grafana]
        JAEG_ST[Jaeger]
    end

    ING_POD --> KAFKA_STRIMZI
    HPA --> ING_POD
    ING_POD -.-> PROM_ST
    GRAF_ST --> PROM_ST
```

## Congestion & routing decision

```mermaid
flowchart LR
    LH[latency-events] --> AGG[Windowed p95 per bank]
    BH[bank-health] --> AGG
    AGG --> SCORE[Congestion score]
    SCORE --> ROUTE[Intelligent Routing]
    AI[AI Prediction] --> ROUTE
    ROUTE -->|reroute| BANK[Bank Simulator]
    ROUTE -->|congestion-events| KAFKA[(Kafka)]
```
