package kafkax

import (
	"context"
	"encoding/json"
	"time"

	kgo "github.com/segmentio/kafka-go"
	"github.com/upi-platform/shared/events"
	"github.com/upi-platform/shared/kafka"
	"github.com/upi-platform/shared/models"
)

type Producer struct {
	writers map[string]*kgo.Writer
}

func NewProducer(brokers []string) *Producer {
	p := &Producer{writers: make(map[string]*kgo.Writer)}
	for _, topic := range []string{
		kafka.TopicUpiTransactions,
		kafka.TopicValidatedTransactions,
		kafka.TopicFailedTransactions,
		kafka.TopicRetryTransactions,
		kafka.TopicAnalyticsEvents,
		kafka.TopicLatencyEvents,
		kafka.TopicBankHealth,
		kafka.TopicDeadLetterEvents,
	} {
		p.writers[topic] = &kgo.Writer{
			Addr:         kgo.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kgo.Hash{},
			RequiredAcks: kgo.RequireOne,
			Async:        false,
		}
	}
	return p
}

func (p *Producer) Publish(ctx context.Context, topic, key, correlationID, idempotencyKey string, env models.EventEnvelope) error {
	body, err := json.Marshal(env)
	if err != nil {
		return err
	}
	msg := kgo.Message{
		Key:   []byte(key),
		Value: body,
		Time:  time.Now().UTC(),
		Headers: []kgo.Header{
			{Key: kafka.HeaderCorrelationID, Value: []byte(correlationID)},
			{Key: kafka.HeaderEventVersion, Value: []byte(env.EventVersion)},
			{Key: kafka.HeaderContentType, Value: []byte("application/json")},
		},
	}
	if idempotencyKey != "" {
		msg.Headers = append(msg.Headers, kgo.Header{Key: kafka.HeaderIdempotencyKey, Value: []byte(idempotencyKey)})
	}
	return p.writers[topic].WriteMessages(ctx, msg)
}

func (p *Producer) PublishEvent(ctx context.Context, topic, key, correlationID, idempotencyKey, eventType, source string, payload any) error {
	env := events.BuildEnvelope(eventType, source, correlationID, payload)
	return p.Publish(ctx, topic, key, correlationID, idempotencyKey, env)
}

func (p *Producer) Close() error {
	var first error
	for _, w := range p.writers {
		if err := w.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}

type Consumer struct {
	reader *kgo.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kgo.NewReader(kgo.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       1,
			MaxBytes:       10e6,
			CommitInterval: time.Second,
			StartOffset:    kgo.LastOffset,
		}),
	}
}

func (c *Consumer) Read(ctx context.Context) (kgo.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func ParseEnvelope(msg kgo.Message) (models.EventEnvelope, error) {
	return events.UnmarshalEnvelope(msg.Value)
}
