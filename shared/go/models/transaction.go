package models

import "time"

// TxType represents UPI transaction categories.
type TxType string

const (
	TxTypeP2P      TxType = "P2P"
	TxTypeP2M      TxType = "P2M"
	TxTypeMandate  TxType = "MANDATE"
	TxTypeCollect  TxType = "COLLECT"
)

// UpiTransaction is the canonical transaction payload (schema: transaction.json).
type UpiTransaction struct {
	TransactionID     string            `json:"transaction_id"`
	IdempotencyKey    string            `json:"idempotency_key"`
	AmountPaise       int64             `json:"amount_paise"`
	Currency          string            `json:"currency"`
	PayerVPA          string            `json:"payer_vpa"`
	PayeeVPA          string            `json:"payee_vpa"`
	BankCode          string            `json:"bank_code"`
	MerchantID        string            `json:"merchant_id"`
	MerchantCategory  string            `json:"merchant_category,omitempty"`
	TxnType           TxType            `json:"txn_type"`
	Note              string            `json:"note,omitempty"`
	DeviceFingerprint string            `json:"device_fingerprint"`
	GeoRegion         string            `json:"geo_region,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// TransactionStatus lifecycle states.
type TransactionStatus string

const (
	StatusReceived   TransactionStatus = "RECEIVED"
	StatusValidated  TransactionStatus = "VALIDATED"
	StatusProcessing TransactionStatus = "PROCESSING"
	StatusSuccess    TransactionStatus = "SUCCESS"
	StatusFailed     TransactionStatus = "FAILED"
	StatusRetrying   TransactionStatus = "RETRYING"
	StatusDLQ        TransactionStatus = "DLQ"
)
