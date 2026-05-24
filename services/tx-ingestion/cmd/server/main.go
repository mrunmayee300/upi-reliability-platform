package main

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/idempotency"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/kafkax"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/metrics"
	"github.com/upi-platform/shared/models"
	"github.com/upi-platform/shared/httpx"
	"github.com/upi-platform/shared/store"
)

const serviceName = "tx-ingestion"

var vpaPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9]+$`)

func main() {
	cfg := config.LoadBase(serviceName, "8081")
	log := logging.New(serviceName)

	ctx := context.Background()
	pg, err := store.NewPostgres(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Error("postgres connect failed", slog.Any("error", err))
		panic(err)
	}
	defer pg.Close()

	idem := idempotency.NewStore(cfg.RedisAddr, cfg.RedisPassword, 24*time.Hour)
	if err := idem.Ping(ctx); err != nil {
		log.Error("redis connect failed", slog.Any("error", err))
		panic(err)
	}

	producer := kafkax.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	r := httpx.NewRouter()
	httpx.RegisterHealth(r, nil, func() error {
		if err := pg.Ping(ctx); err != nil {
			return err
		}
		return idem.Ping(ctx)
	})

	r.POST("/v1/transactions", func(c *gin.Context) {
		metrics.HTTPRequests.WithLabelValues(serviceName, "POST", "/v1/transactions", "202").Inc()

		idemKey := c.GetHeader("Idempotency-Key")
		if idemKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error_code": "MISSING_IDEMPOTENCY_KEY"})
			return
		}
		correlationID := httpx.CorrelationID(c)

		var txn models.UpiTransaction
		if err := c.ShouldBindJSON(&txn); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_PAYLOAD", "message": err.Error()})
			return
		}

		if err := validateTransaction(&txn); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error_code": "VALIDATION_ERROR", "message": err.Error()})
			return
		}

		if txn.TransactionID == "" {
			txn.TransactionID = "TXN-" + uuid.New().String()
		}
		txn.IdempotencyKey = idemKey
		if txn.Currency == "" {
			txn.Currency = "INR"
		}
		if txn.CreatedAt.IsZero() {
			txn.CreatedAt = time.Now().UTC()
		}

		acquired, dup, err := idem.TryAcquire(ctx, idemKey, txn.TransactionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "IDEMPOTENCY_ERROR"})
			return
		}
		if !acquired {
			c.JSON(http.StatusConflict, gin.H{
				"error_code":     "DUPLICATE_TXN",
				"transaction_id": dup.TransactionID,
				"correlation_id": correlationID,
			})
			return
		}

		if err := pg.InsertTransaction(ctx, txn, models.StatusValidated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": "DB_ERROR"})
			return
		}

		_ = producer.PublishEvent(ctx, kafka.TopicUpiTransactions, txn.PayerVPA, correlationID, idemKey,
			models.EventTransactionReceived, serviceName, txn)
		_ = producer.PublishEvent(ctx, kafka.TopicValidatedTransactions, txn.TransactionID, correlationID, idemKey,
			models.EventTransactionValidated, serviceName, txn)

		metrics.KafkaMessages.WithLabelValues(serviceName, kafka.TopicValidatedTransactions, "produce").Inc()
		metrics.TxProcessed.WithLabelValues(serviceName, "validated").Inc()

		c.JSON(http.StatusAccepted, gin.H{
			"transaction_id": txn.TransactionID,
			"status":         "ACCEPTED",
			"correlation_id": correlationID,
		})
	})

	r.GET("/v1/transactions/:transactionId", func(c *gin.Context) {
		txnID := c.Param("transactionId")
		status, updatedAt, err := pg.GetTransaction(ctx, txnID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error_code": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"transaction_id": txnID,
			"status":         status,
			"updated_at":     updatedAt,
		})
	})

	log.Info("starting", slog.String("port", cfg.HTTPPort))
	if err := httpx.Run(r, cfg.HTTPPort); err != nil {
		log.Error("server stopped", slog.Any("error", err))
	}
}

func validateTransaction(txn *models.UpiTransaction) error {
	if txn.AmountPaise < 100 {
		return errString("amount_paise must be >= 100")
	}
	if !vpaPattern.MatchString(txn.PayerVPA) || !vpaPattern.MatchString(txn.PayeeVPA) {
		return errString("invalid VPA format")
	}
	if strings.TrimSpace(txn.BankCode) == "" {
		return errString("bank_code required")
	}
	if strings.TrimSpace(txn.MerchantID) == "" {
		return errString("merchant_id required")
	}
	return nil
}

type errString string

func (e errString) Error() string { return string(e) }
