#!/usr/bin/env bash
set -euo pipefail

if ! docker info >/dev/null 2>&1; then
  echo "ERROR: Docker is not running. Start Docker Desktop and retry."
  exit 1
fi

if ! docker inspect --format '{{.State.Running}}' upi-kafka 2>/dev/null | grep -q true; then
  echo "ERROR: Container 'upi-kafka' is not running. Run: docker compose up -d kafka"
  exit 1
fi

BROKER="${KAFKA_BROKER:-localhost:9092}"

create_topic() {
  local name="$1"
  local partitions="$2"
  local retention_ms="$3"
  local cleanup_policy="${4:-delete}"

  echo "Ensuring topic: ${name}"
  docker exec upi-kafka /opt/kafka/bin/kafka-topics.sh \
    --bootstrap-server "${BROKER}" \
    --create \
    --if-not-exists \
    --topic "${name}" \
    --partitions "${partitions}" \
    --replication-factor 1 \
    --config retention.ms="${retention_ms}" \
    --config cleanup.policy="${cleanup_policy}" >/dev/null
}

create_topic "upi-transactions" 24 259200000
create_topic "validated-transactions" 24 172800000
create_topic "failed-transactions" 12 259200000
create_topic "retry-transactions" 12 172800000
create_topic "fraud-alerts" 6 604800000
create_topic "latency-events" 12 86400000
create_topic "bank-health" 6 86400000 compact
create_topic "congestion-events" 6 172800000
create_topic "analytics-events" 12 172800000
create_topic "dead-letter-events" 6 2592000000

echo "Kafka topic bootstrap complete."
