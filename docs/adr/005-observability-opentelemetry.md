# ADR-005: Unified Observability with OpenTelemetry

## Status

Accepted

## Context

Distributed payment flows require correlated logs, metrics, and traces across HTTP and Kafka.

## Decision

OpenTelemetry SDK in all Go services; OTLP export to collector → Jaeger (traces) + Prometheus (metrics via OTel Collector Prometheus exporter).

Structured JSON logs with `correlation_id`, `trace_id`, `service.name`.

## Consequences

**Positive:** CNCF-standard, portable to Grafana Cloud/Datadog later.  
**Negative:** Collector deployment overhead (included in docker-compose Phase 2).
