# ADR-002: Go for Core Microservices

## Status

Accepted

## Context

Payment paths need low latency, predictable memory, and strong concurrency for ingestion and bank simulation.

## Decision

Implement 10 core services in Go 1.22+ with Gin, `segmentio/kafka-go`, OpenTelemetry SDK.

## Consequences

**Positive:** Single binary deploys, excellent throughput, aligns with cloud-native tooling.  
**Negative:** ML stays in Python (second runtime).

## Alternatives considered

- Java/Spring — heavier footprint for demo cluster
- Rust — higher dev cost for broad microservice surface
