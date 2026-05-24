package idempotency

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	ttl    time.Duration
}

type Record struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

func NewStore(addr, password string, ttl time.Duration) *Store {
	return &Store{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
		}),
		ttl: ttl,
	}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *Store) TryAcquire(ctx context.Context, key, transactionID string) (bool, *Record, error) {
	redisKey := fmt.Sprintf("idempotency:%s", key)
	ok, err := s.client.SetNX(ctx, redisKey, transactionID, s.ttl).Result()
	if err != nil {
		return false, nil, err
	}
	if ok {
		return true, nil, nil
	}
	val, err := s.client.Get(ctx, redisKey).Result()
	if err != nil {
		return false, nil, err
	}
	return false, &Record{TransactionID: val, Status: "DUPLICATE"}, nil
}

func (s *Store) SaveRecord(ctx context.Context, key string, rec Record) error {
	redisKey := fmt.Sprintf("idempotency:%s", key)
	body, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, redisKey, body, s.ttl).Err()
}
