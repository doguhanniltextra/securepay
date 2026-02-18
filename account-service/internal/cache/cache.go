package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// BalanceEntry represents a cached balance value.
type BalanceEntry struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// Cache interface defines operations for balance caching.
type Cache interface {
	GetBalance(ctx context.Context, accountID string) (*BalanceEntry, error)
	SetBalance(ctx context.Context, accountID string, entry *BalanceEntry) error
	DeleteBalance(ctx context.Context, accountID string) error
}

type redisCache struct {
	client *redis.Client
}

const (
	balanceKeyPrefix = "balance:"
	balanceTTL       = 60 * time.Second // Convention: 60 seconds
)

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(addr, password string) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &redisCache{client: client}
}

// GetBalance retrieves a cached balance for the given account ID.
// Returns nil, nil if key does not exist (cache miss).
func (r *redisCache) GetBalance(ctx context.Context, accountID string) (*BalanceEntry, error) {
	key := balanceKeyPrefix + accountID
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var entry BalanceEntry
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balance entry: %w", err)
	}

	return &entry, nil
}

// SetBalance stores a balance entry in cache with TTL of 60 seconds.
func (r *redisCache) SetBalance(ctx context.Context, accountID string, entry *BalanceEntry) error {
	key := balanceKeyPrefix + accountID
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal balance entry: %w", err)
	}

	if err := r.client.Set(ctx, key, string(data), balanceTTL).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

// DeleteBalance removes a cached balance entry.
func (r *redisCache) DeleteBalance(ctx context.Context, accountID string) error {
	key := balanceKeyPrefix + accountID
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}
