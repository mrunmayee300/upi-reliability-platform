# Scripts

| Script | Phase | Purpose |
|--------|-------|---------|
| `kafka/create-topics.sh` | 2 | Idempotent Kafka topic creation (bash) |
| `kafka/create-topics.ps1` | 2 | Idempotent Kafka topic creation (PowerShell) |
| `dev/bootstrap.sh` | 2 | Start infra + create topics (bash) |
| `dev/bootstrap.ps1` | 2 | Start infra + create topics (PowerShell) |
| `dev/build-all.ps1` | 3 | Build all Go microservices |
| `dev/run-services.ps1` | 3 | Start services in background (logs/pids) |
| `dev/stop-services.ps1` | 3 | Stop all background services |
| `dev/start-generator.ps1` | 3 | Start load generator (checks port 8082) |
| `dev/run-dashboard.ps1` | 4 | Start Next.js ops dashboard (:3000) |
| `dev/verify-services.ps1` | 3 | Health-check all service ports |
| `k8s/create-kind.ps1` | 5 | Create kind cluster |
| `k8s/build-images.ps1` | 5 | Build & load Docker images into kind |
| `k8s/deploy-helm.ps1` | 5 | Helm install/upgrade |
| `dev/seed-banks.sh` | 3 | Seed bank simulation configs |
| `chaos/kafka-outage.sh` | 6 | Chaos test: broker outage |
| `benchmark/run-k6.sh` | 6 | Run load benchmark suite |
