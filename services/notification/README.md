# Notification / Alert Service

**Port:** 8090 | **Phase 3**

## Responsibilities

- Consume `fraud-alerts`, `dead-letter-events`, `congestion-events`
- Incident routing (webhook simulation, log sink)
- Alert deduplication and severity thresholds

## Kafka

- **Consumes:** `fraud-alerts`, `dead-letter-events`, `congestion-events`
