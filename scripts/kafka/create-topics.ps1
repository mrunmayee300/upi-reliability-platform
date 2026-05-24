# Idempotent Kafka topic creation (Windows / PowerShell)
$ErrorActionPreference = "Stop"

function Test-DockerRunning {
    $prev = $ErrorActionPreference
    $ErrorActionPreference = "SilentlyContinue"
    docker info *> $null
    $ok = ($LASTEXITCODE -eq 0)
    $ErrorActionPreference = $prev
    return $ok
}

function Test-KafkaContainer {
    $state = docker inspect --format '{{.State.Running}}' upi-kafka 2>$null
    return $state -eq "true"
}

function New-KafkaTopic {
    param(
        [string]$Name,
        [int]$Partitions,
        [long]$RetentionMs,
        [string]$CleanupPolicy = "delete"
    )

    Write-Host "Ensuring topic: $Name"
    docker exec upi-kafka /opt/kafka/bin/kafka-topics.sh `
        --bootstrap-server localhost:9092 `
        --create `
        --if-not-exists `
        --topic $Name `
        --partitions $Partitions `
        --replication-factor 1 `
        --config "retention.ms=$RetentionMs" `
        --config "cleanup.policy=$CleanupPolicy" | Out-Null

    if ($LASTEXITCODE -ne 0) {
        throw "Failed to create topic: $Name"
    }
}

if (-not (Test-DockerRunning)) {
    Write-Error @"
Docker is not running.

Start Docker Desktop and wait until it shows 'Engine running', then run:
  .\scripts\dev\bootstrap.ps1
"@
    exit 1
}

if (-not (Test-KafkaContainer)) {
    Write-Error "Container 'upi-kafka' is not running. Run: docker compose up -d kafka"
    exit 1
}

New-KafkaTopic "upi-transactions" 24 259200000
New-KafkaTopic "validated-transactions" 24 172800000
New-KafkaTopic "failed-transactions" 12 259200000
New-KafkaTopic "retry-transactions" 12 172800000
New-KafkaTopic "fraud-alerts" 6 604800000
New-KafkaTopic "latency-events" 12 86400000
New-KafkaTopic "bank-health" 6 86400000 "compact"
New-KafkaTopic "congestion-events" 6 172800000
New-KafkaTopic "analytics-events" 12 172800000
New-KafkaTopic "dead-letter-events" 6 2592000000

Write-Host "Kafka topic bootstrap complete." -ForegroundColor Green
