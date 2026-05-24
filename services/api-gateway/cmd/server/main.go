package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/upi-platform/shared/config"
	"github.com/upi-platform/shared/logging"
	"github.com/upi-platform/shared/httpx"
	"github.com/upi-platform/shared/metrics"
)

const serviceName = "api-gateway"

func main() {
	cfg := config.LoadBase(serviceName, "8080")
	log := logging.New(serviceName)

	ingestionURL := getenv("INGESTION_URL", "http://localhost:8081")
	analyticsURL := getenv("ANALYTICS_URL", "http://localhost:8089")
	apiKeys := parseAPIKeys(getenv("API_KEYS", "dev-api-key-001"))
	rateLimit := config.GetInt("RATE_LIMIT_RPS", 1000)

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword})
	client := &http.Client{Timeout: 15 * time.Second}

	r := httpx.NewRouter()
	httpx.RegisterHealth(r, nil, func() error {
		return rdb.Ping(context.Background()).Err()
	})

	auth := func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !apiKeys[token] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHORIZED"})
			return
		}
		c.Next()
	}

	rateLimitMW := func(c *gin.Context) {
		key := "ratelimit:" + c.GetHeader("Authorization")
		ctx := context.Background()
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			_ = rdb.Expire(ctx, key, time.Second).Err()
		}
		if count > int64(rateLimit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error_code": "RATE_LIMITED"})
			return
		}
		c.Next()
	}

	r.POST("/v1/transactions", auth, rateLimitMW, func(c *gin.Context) {
		proxyPOST(c, client, ingestionURL+"/v1/transactions")
		metrics.HTTPRequests.WithLabelValues(serviceName, "POST", "/v1/transactions", "202").Inc()
	})

	r.GET("/v1/metrics/summary", auth, func(c *gin.Context) {
		proxyGET(c, client, analyticsURL+"/v1/metrics/summary")
	})

	log.Info("api gateway running", slog.String("port", cfg.HTTPPort))
	_ = httpx.Run(r, cfg.HTTPPort)
}

func proxyPOST(c *gin.Context, client *http.Client, url string) {
	body, _ := io.ReadAll(c.Request.Body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header = c.Request.Header.Clone()
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error_code": "UPSTREAM_ERROR"})
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func proxyGET(c *gin.Context, client *http.Client, url string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error_code": "UPSTREAM_ERROR"})
		return
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func parseAPIKeys(raw string) map[string]bool {
	m := make(map[string]bool)
	for _, k := range strings.Split(raw, ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			m[k] = true
		}
	}
	return m
}
