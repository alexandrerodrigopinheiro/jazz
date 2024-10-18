package cache

import (
	"context"
	"encoding/json"
	"time"

	"jazz/backend/configs"
	"jazz/backend/pkg/logger"

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
		logger.Logger.Warnw("REDIS_URL not set in configuration. Falling back to default cache.")
		return nil
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.Logger.Errorw("Failed to parse REDIS_URL", "error", err)
		return nil
	}

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Logger.Errorw("Failed to connect to Redis", "error", err)
		return nil
	}

	logger.Logger.Infow("Redis cache initialized successfully", "url", redisURL)
	return &RedisCache{client: client}
}

// Set stores a value in Redis.
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		logger.Logger.Errorw("Error serializing value", "key", key, "error", err)
		return err
	}

	err = r.client.Set(context.Background(), key, valueBytes, expiration).Err()
	if err != nil {
		logger.Logger.Errorw("Failed to set value in Redis", "key", key, "error", err)
	}
	logger.Logger.Infow("Value set in Redis cache", "key", key, "expiration", expiration)
	return err
}

// Remember stores a value in Redis using a callback if the value does not already exist.
func (r *RedisCache) Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	logger.Logger.Infow("Attempting to retrieve key from Redis cache", "key", key)
	var result interface{}

	// Check if value is already in the cache
	value, err := r.Get(key)
	if err == nil && value != nil {
		if valStr, ok := value.(string); ok {
			err = json.Unmarshal([]byte(valStr), &result)
			if err == nil {
				logger.Logger.Infow("Cache hit in Redis", "key", key)
				return result, nil
			}
		}
	}

	// If value is not cached, execute the callback
	logger.Logger.Infow("Cache miss, executing callback", "key", key)
	result, err = callback()
	if err != nil {
		logger.Logger.Errorw("Callback execution failed", "key", key, "error", err)
		return nil, err
	}

	// Cache the value
	if err := r.Set(key, result, expiration); err != nil {
		logger.Logger.Errorw("Failed to set value in Redis after callback", "key", key, "error", err)
		return result, err
	}

	logger.Logger.Infow("Value cached in Redis after callback execution", "key", key)
	return result, nil
}

// Forget removes a value from Redis.
func (r *RedisCache) Forget(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		logger.Logger.Errorw("Failed to delete key from Redis", "key", key, "error", err)
		return err
	}
	logger.Logger.Infow("Cache entry removed from Redis", "key", key)
	return nil
}

// Get retrieves a value from Redis.
func (r *RedisCache) Get(key string) (interface{}, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		logger.Logger.Warnw("Cache miss in Redis", "key", key)
		return nil, nil
	}
	if err != nil {
		logger.Logger.Errorw("Failed to get value from Redis", "key", key, "error", err)
		return nil, err
	}

	// Try to unmarshal the value as JSON
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err == nil {
		logger.Logger.Infow("Cache hit in Redis", "key", key)
		return result, nil
	}

	logger.Logger.Infow("Returning raw value from Redis cache", "key", key)
	return val, nil
}
