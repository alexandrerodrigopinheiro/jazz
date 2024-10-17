package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"jazz/backend/configs"

	"github.com/go-redis/redis/v8"
)

// RedisCache Implementation of Cache using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache initializes a new Redis Cache.
func NewRedisCache() *RedisCache {
	cacheConfig := configs.GetCacheConfig()
	storeConfig := cacheConfig["stores"].(map[string]interface{})
	redisConfig := storeConfig["redis"].(map[string]interface{})

	redisURL, ok := redisConfig["url"].(string)
	if !ok || redisURL == "" {
		fmt.Println("REDIS_URL not set in configuration. Falling back to default cache.")
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Println("Failed to parse REDIS_URL:", err)
		return nil
	}

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return nil
	}

	return &RedisCache{client: client}
}

// Set stores a value in Redis.
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error serializing value for key %s: %w", key, err)
	}
	return r.client.Set(context.Background(), key, valueBytes, expiration).Err()
}

// Remember stores a value in Redis using a callback if the value does not already exist.
func (r *RedisCache) Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	var result interface{}

	// Check if value is already in the cache
	value, err := r.Get(key)
	if err == nil && value != nil {
		if valStr, ok := value.(string); ok {
			err = json.Unmarshal([]byte(valStr), &result)
			if err == nil {
				return result, nil
			}
		}
	}

	// If value is not cached, execute the callback
	result, err = callback()
	if err != nil {
		return nil, fmt.Errorf("callback execution failed for key %s: %w", key, err)
	}

	// Cache the value
	if err := r.Set(key, result, expiration); err != nil {
		return result, err
	}

	return result, nil
}

// Forget removes a value from Redis.
func (r *RedisCache) Forget(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

// Get retrieves a value from Redis.
func (r *RedisCache) Get(key string) (interface{}, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
	}

	// Try to unmarshal the value as JSON
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err == nil {
		return result, nil
	}

	return val, nil
}
