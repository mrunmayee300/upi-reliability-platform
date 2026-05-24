# ADR-003: Idempotency via Redis + PostgreSQL

## Status

Accepted

## Context

UPI clients and retry orchestrators may duplicate requests. Duplicate settlement must be prevented.

## Decision

- **Hot path:** Redis `SET key NX EX` with key `idempotency:{key}`
- **Audit:** Postgres `idempotency_records` with response snapshot

## Consequences

**Positive:** Sub-ms dedupe; durable audit for disputes.  
**Negative:** Redis failure mode requires fail-closed or degraded policy (config: `IDEMPOTENCY_FAIL_CLOSED=true`).
