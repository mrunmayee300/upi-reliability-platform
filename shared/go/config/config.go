package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// TrimSpace applied in getEnv — avoids Windows cmd.exe trailing spaces in env vars.

// Base holds common service configuration from environment.
type Base struct {
	ServiceName   string
	HTTPPort      string
	KafkaBrokers  []string
	PostgresDSN   string
	RedisAddr     string
	RedisPassword string
	OTelEndpoint  string
}

func LoadBase(serviceName, defaultPort string) Base {
	return Base{
		ServiceName:   serviceName,
		HTTPPort:      getEnv("HTTP_PORT", defaultPort),
		KafkaBrokers:  strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		PostgresDSN:   getEnv("POSTGRES_DSN", "postgres://upi:upi_dev_password@localhost:5433/upi_platform?sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		OTelEndpoint:  getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318"),
	}
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func GetFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
