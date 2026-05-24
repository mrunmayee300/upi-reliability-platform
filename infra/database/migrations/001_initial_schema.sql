-- UPI Platform — Initial PostgreSQL schema (Phase 2 apply)
-- Requires: PostgreSQL 15+

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Idempotency audit trail
CREATE TABLE idempotency_records (
    idempotency_key   VARCHAR(128) PRIMARY KEY,
    transaction_id    VARCHAR(64) NOT NULL,
    response_status   VARCHAR(32) NOT NULL,
    response_body     JSONB,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at        TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_idempotency_expires ON idempotency_records (expires_at);

-- Transaction master (durable read model)
CREATE TABLE transactions (
    transaction_id    VARCHAR(64) PRIMARY KEY,
    idempotency_key   VARCHAR(128) NOT NULL UNIQUE,
    status            VARCHAR(32) NOT NULL,
    amount_paise      BIGINT NOT NULL,
    payer_vpa         VARCHAR(255) NOT NULL,
    payee_vpa         VARCHAR(255) NOT NULL,
    bank_code         VARCHAR(16) NOT NULL,
    merchant_id       VARCHAR(64) NOT NULL,
    txn_type          VARCHAR(16) NOT NULL,
    payload           JSONB NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_transactions_status ON transactions (status);
CREATE INDEX idx_transactions_bank ON transactions (bank_code);
CREATE INDEX idx_transactions_created ON transactions (created_at DESC);

-- Retry attempts
CREATE TABLE retry_attempts (
    id                BIGSERIAL PRIMARY KEY,
    transaction_id    VARCHAR(64) NOT NULL REFERENCES transactions(transaction_id),
    attempt_number    INT NOT NULL,
    scheduled_at      TIMESTAMPTZ NOT NULL,
    executed_at       TIMESTAMPTZ,
    failure_code      VARCHAR(64),
    status            VARCHAR(32) NOT NULL,
    UNIQUE (transaction_id, attempt_number)
);
CREATE INDEX idx_retry_scheduled ON retry_attempts (scheduled_at) WHERE status = 'PENDING';

-- Dead letter queue (mirror of Kafka for query API)
CREATE TABLE dead_letter_queue (
    event_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id    VARCHAR(64) NOT NULL,
    failure_history   JSONB NOT NULL,
    final_failure_code VARCHAR(64) NOT NULL,
    correlation_id    UUID NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Transaction audit log
CREATE TABLE transaction_audit (
    id                BIGSERIAL PRIMARY KEY,
    transaction_id    VARCHAR(64) NOT NULL,
    correlation_id    UUID NOT NULL,
    actor_service     VARCHAR(64) NOT NULL,
    event_type        VARCHAR(64) NOT NULL,
    payload_hash      VARCHAR(64),
    details           JSONB,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_txn ON transaction_audit (transaction_id, created_at DESC);

-- Analytics rollups (1-minute windows)
CREATE TABLE metrics_rollup_1m (
    window_start      TIMESTAMPTZ NOT NULL,
    metric_name       VARCHAR(64) NOT NULL,
    bank_code         VARCHAR(16),
    value             DOUBLE PRECISION NOT NULL,
    sample_count      BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (window_start, metric_name, bank_code)
);

-- Fraud alerts persistence
CREATE TABLE fraud_alerts (
    alert_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id    VARCHAR(64) NOT NULL,
    anomaly_score     DOUBLE PRECISION NOT NULL,
    severity          VARCHAR(16) NOT NULL,
    signals           JSONB,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_fraud_created ON fraud_alerts (created_at DESC);
