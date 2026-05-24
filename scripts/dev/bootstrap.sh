#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

if ! docker info >/dev/null 2>&1; then
  echo ""
  echo "ERROR: Docker is not running."
  echo ""
  echo "Start Docker Desktop (or your Docker daemon), then run:"
  echo "  ./scripts/dev/bootstrap.sh"
  echo ""
  exit 1
fi

echo "Starting infra stack..."
docker compose up -d kafka postgres redis jaeger otel-collector prometheus grafana kafka-ui

echo "Waiting for Kafka health..."
healthy=0
for _ in $(seq 1 60); do
  status="$(docker inspect --format '{{.State.Health.Status}}' upi-kafka 2>/dev/null || true)"
  if [[ "$status" == "healthy" ]]; then
    healthy=1
    break
  fi
  sleep 2
done

if [[ "$healthy" -ne 1 ]]; then
  echo "ERROR: Kafka did not become healthy. Check: docker compose logs kafka"
  exit 1
fi

echo "Creating Kafka topics..."
bash ./scripts/kafka/create-topics.sh

echo ""
echo "Bootstrap complete."
echo "Kafka UI:   http://localhost:8088"
echo "Prometheus: http://localhost:9090"
echo "Grafana:    http://localhost:3001 (admin/admin)"
echo "Jaeger:     http://localhost:16686"
