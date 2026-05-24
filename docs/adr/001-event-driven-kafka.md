# ADR-001: Event-Driven Architecture with Apache Kafka

## Status

Accepted

## Context

The platform must handle high-volume, ordered payment lifecycle events with replay capability for recovery and analytics.

## Decision

Use Kafka as the central event backbone. PostgreSQL stores durable aggregates; Redis handles hot idempotency and routing caches.

## Consequences

**Positive:** Decoupled services, replay, industry-aligned payment log pattern.  
**Negative:** Operational complexity, eventual consistency in dashboards, requires lag monitoring.

## Alternatives considered

- RabbitMQ — weaker log retention/replay story at scale
- NATS JetStream — viable; less ecosystem for fintech samples
