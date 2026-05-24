package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/models"
	"github.com/upi-platform/shared/httpx"
)

const serviceName = "tx-generator"

var (
	running        atomic.Bool
	totalGenerated atomic.Int64
	totalAccepted  atomic.Int64
	totalRejected  atomic.Int64
	stopCh         chan struct{}
	genMu          sync.Mutex
)

var banks = []string{"HDFC", "ICICI", "SBI", "AXIS", "KOTAK", "PNB", "BOB", "YES"}
var vpas = []string{"paytm", "ybl", "okaxis", "ibl"}

func main() {
	cfg := config.LoadBase(serviceName, "8082")
	log := logging.New(serviceName)
	ingestionURL := getenv("INGESTION_URL", "http://localhost:8081")

	r := httpx.NewRouter()
	httpx.RegisterHealth(r, nil, nil)

	r.POST("/v1/generator/start", func(c *gin.Context) {
		var req struct {
			TPS             int     `json:"tps"`
			FailureRate     float64 `json:"failure_rate"`
			DurationSeconds int     `json:"duration_seconds"`
		}
		_ = c.ShouldBindJSON(&req)
		if req.TPS <= 0 {
			req.TPS = 100
		}
		if req.TPS > 5000 {
			req.TPS = 5000
		}

		genMu.Lock()
		if running.Load() {
			genMu.Unlock()
			c.JSON(http.StatusOK, gin.H{"message": "already running"})
			return
		}
		stopCh = make(chan struct{})
		running.Store(true)
		genMu.Unlock()

		go runGenerator(log, ingestionURL, req.TPS, req.DurationSeconds)
		c.JSON(http.StatusOK, gin.H{"running": true, "tps": req.TPS})
	})

	r.POST("/v1/generator/stop", func(c *gin.Context) {
		genMu.Lock()
		if running.Load() {
			close(stopCh)
			running.Store(false)
		}
		genMu.Unlock()
		c.JSON(http.StatusOK, gin.H{"running": false})
	})

	r.GET("/v1/generator/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"running":         running.Load(),
			"total_generated": totalGenerated.Load(),
			"total_accepted":  totalAccepted.Load(),
			"total_rejected":  totalRejected.Load(),
		})
	})

	log.Info("tx generator running", slog.String("port", cfg.HTTPPort))
	_ = httpx.Run(r, cfg.HTTPPort)
}

func runGenerator(log *slog.Logger, ingestionURL string, tps, durationSec int) {
	interval := time.Second / time.Duration(tps)
	if interval <= 0 {
		interval = time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var deadline time.Time
	if durationSec > 0 {
		deadline = time.Now().Add(time.Duration(durationSec) * time.Second)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			if !deadline.IsZero() && time.Now().After(deadline) {
				running.Store(false)
				return
			}
			txn := randomTxn()
			totalGenerated.Add(1)
			if sendTxn(client, ingestionURL, txn) {
				totalAccepted.Add(1)
			} else {
				totalRejected.Add(1)
			}
		}
	}
}

func randomTxn() models.UpiTransaction {
	bank := banks[rand.Intn(len(banks))]
	payer := "user" + uuid.New().String()[:6] + "@" + vpas[rand.Intn(len(vpas))]
	payee := "merchant" + uuid.New().String()[:4] + "@" + vpas[rand.Intn(len(vpas))]
	return models.UpiTransaction{
		TransactionID:     "TXN-" + uuid.New().String(),
		IdempotencyKey:    uuid.New().String(),
		AmountPaise:       int64(100 + rand.Intn(500000)),
		Currency:          "INR",
		PayerVPA:          payer,
		PayeeVPA:          payee,
		BankCode:          bank,
		MerchantID:        "MRC-" + uuid.New().String()[:8],
		TxnType:           models.TxTypeP2M,
		DeviceFingerprint: uuid.New().String(),
		CreatedAt:         time.Now().UTC(),
	}
}

func sendTxn(client *http.Client, url string, txn models.UpiTransaction) bool {
	body, _ := json.Marshal(txn)
	req, _ := http.NewRequest(http.MethodPost, url+"/v1/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", txn.IdempotencyKey)
	req.Header.Set("X-Correlation-Id", uuid.New().String())
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusAccepted
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
