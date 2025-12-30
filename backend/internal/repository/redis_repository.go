package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// CacheTTL is the default time-to-live for cached URLs
	CacheTTL = 24 * time.Hour
)

// RedisRepository implements domain.CacheRepository (Adapter)
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new Redis repository
func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

// Get retrieves a value from cache
func (r *RedisRepository) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key doesn't exist, not an error
	}
	return val, err
}

// Set stores a value in cache with default TTL
func (r *RedisRepository) Set(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return r.client.Set(ctx, key, value, CacheTTL).Err()
}

// SetWithTTL stores a value in cache with custom TTL
func (r *RedisRepository) SetWithTTL(key string, value string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a value from cache
func (r *RedisRepository) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (r *RedisRepository) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}
