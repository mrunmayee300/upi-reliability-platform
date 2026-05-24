package otel

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const correlationIDKey ctxKey = "correlation_id"

// WithCorrelationID attaches correlation ID to context.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

// CorrelationIDFromContext returns correlation ID or generates one.
func CorrelationIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(correlationIDKey).(string); ok && v != "" {
		return v
	}
	return uuid.New().String()
}
