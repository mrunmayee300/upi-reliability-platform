package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/upi-platform/shared/models"
)

// BuildEnvelope creates a canonical Kafka event envelope.
func BuildEnvelope(eventType, source string, correlationID string, payload any) models.EventEnvelope {
	return models.EventEnvelope{
		EventID:       uuid.New().String(),
		EventType:     eventType,
		EventVersion:  "1.0",
		OccurredAt:    time.Now().UTC(),
		CorrelationID: correlationID,
		SourceService: source,
		Payload:       payload,
	}
}

func MarshalEnvelope(env models.EventEnvelope) ([]byte, error) {
	return json.Marshal(env)
}

func UnmarshalEnvelope(data []byte) (models.EventEnvelope, error) {
	var env models.EventEnvelope
	err := json.Unmarshal(data, &env)
	return env, err
}
