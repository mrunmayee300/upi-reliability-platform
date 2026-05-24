# ADR-004: Failure Processing — Detector + Orchestrator Split

## Status

Accepted

## Context

Failures need classification (Detector) separate from side-effecting retries (Orchestrator).

## Decision

- **Failure Detector:** consumes `failed-transactions`, enriches `failure_cause`, republishes enriched event to same topic (compacted metadata header `enriched=true`) OR to `analytics-events` — Phase 3 implements enrich-in-place with new `event_type`.
- **Retry Orchestrator:** consumes enriched failures, manages backoff and `retry-transactions`.

## Consequences

**Positive:** Single responsibility, independent scaling.  
**Negative:** Careful consumer offset management to avoid double-retry; mitigated by idempotency keys.
