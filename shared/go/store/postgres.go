package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/upi-platform/shared/models"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Postgres) InsertTransaction(ctx context.Context, txn models.UpiTransaction, status models.TransactionStatus) error {
	payload, err := json.Marshal(txn)
	if err != nil {
		return err
	}
	_, err = p.pool.Exec(ctx, `
		INSERT INTO transactions (
			transaction_id, idempotency_key, status, amount_paise, payer_vpa, payee_vpa,
			bank_code, merchant_id, txn_type, payload, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW())
		ON CONFLICT (transaction_id) DO UPDATE SET status = EXCLUDED.status, updated_at = NOW()
	`, txn.TransactionID, txn.IdempotencyKey, status, txn.AmountPaise, txn.PayerVPA, txn.PayeeVPA,
		txn.BankCode, txn.MerchantID, txn.TxnType, payload, txn.CreatedAt)
	return err
}

func (p *Postgres) UpdateStatus(ctx context.Context, transactionID string, status models.TransactionStatus) error {
	_, err := p.pool.Exec(ctx, `UPDATE transactions SET status = $2, updated_at = NOW() WHERE transaction_id = $1`, transactionID, status)
	return err
}

func (p *Postgres) GetTransaction(ctx context.Context, transactionID string) (models.TransactionStatus, time.Time, error) {
	var status string
	var updatedAt time.Time
	err := p.pool.QueryRow(ctx, `SELECT status, updated_at FROM transactions WHERE transaction_id = $1`, transactionID).Scan(&status, &updatedAt)
	return models.TransactionStatus(status), updatedAt, err
}

func (p *Postgres) InsertRetryAttempt(ctx context.Context, transactionID string, attempt int, scheduledAt time.Time) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO retry_attempts (transaction_id, attempt_number, scheduled_at, status)
		VALUES ($1,$2,$3,'PENDING')
		ON CONFLICT (transaction_id, attempt_number) DO NOTHING
	`, transactionID, attempt, scheduledAt)
	return err
}

func (p *Postgres) InsertDLQ(ctx context.Context, transactionID, failureCode, correlationID string, history []byte) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO dead_letter_queue (transaction_id, failure_history, final_failure_code, correlation_id)
		VALUES ($1,$2,$3,$4)
	`, transactionID, history, failureCode, correlationID)
	return err
}
