package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"jazz/backend/configs"
	"jazz/backend/pkg/logger"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedCache Implementation of Cache using Memcached.
type MemcachedCache struct {
	client *memcache.Client
}

// NewMemcachedCache initializes a new Memcached Cache.
func NewMemcachedCache() *MemcachedCache {
	cacheConfig := configs.GetCacheConfig()["stores"].(map[string]interface{})["memcached"].(map[string]interface{})

	host, okHost := cacheConfig["host"].(string)
	port, okPort := cacheConfig["port"].(int)

	if !okHost || host == "" || !okPort {
		logger.Logger.Warn("MEMCACHED_HOST or MEMCACHED_PORT is not set or invalid. Falling back to default cache.")
		return nil
	}

	portStr := fmt.Sprintf("%d", port)
	client := memcache.New(fmt.Sprintf("%s:%s", host, portStr))

	// Testing the connection with Memcached
	testKey := "test_connection"
	testValue := []byte("ping")
	err := client.Set(&memcache.Item{Key: testKey, Value: testValue, Expiration: 1})
	if err != nil {
		logger.Logger.Warnf("Memcached unavailable at %s:%s. Error: %s. Falling back to default cache.", host, portStr, err)
		return nil
	}

	_, err = client.Get(testKey)
	if err != nil {
		logger.Logger.Warnf("Memcached unavailable at %s:%s. Error: %s. Falling back to default cache.", host, portStr, err)
		return nil
	}

	logger.Logger.Infof("Connected to Memcached at %s:%s", host, portStr)
	return &MemcachedCache{client: client}
}

// Set stores a value in Memcached.
func (m *MemcachedCache) Set(key string, value interface{}, expiration time.Duration) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		logger.Logger.Errorw("Failed to serialize value", "key", key, "error", err)
		return err
	}

	item := &memcache.Item{
		Key:        key,
		Value:      valueBytes,
		Expiration: int32(expiration.Seconds()),
	}

	return m.client.Set(item)
}

// Remember stores a value in Memcached using a callback if the value does not already exist.
func (m *MemcachedCache) Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Check if value is already in the cache
	value, err := m.Get(key)
	if err == nil && value != nil {
		var result interface{}
		err = json.Unmarshal([]byte(value.(string)), &result)
		if err == nil {
			return result, nil
		}
	}

	// If value is not cached, execute the callback
	result, err := callback()
	if err != nil {
		logger.Logger.Errorw("Callback execution failed", "key", key, "error", err)
		return result, err
	}

	// Cache the value
	if err := m.Set(key, result, expiration); err != nil {
		logger.Logger.Errorw("Failed to set value in cache after callback", "key", key, "error", err)
		return result, err
	}

	return result, nil
}

// Forget removes a value from Memcached.
func (m *MemcachedCache) Forget(key string) error {
	err := m.client.Delete(key)
	if err != nil && err != memcache.ErrCacheMiss {
		logger.Logger.Errorw("Failed to delete value from Memcached", "key", key, "error", err)
		return err
	}

	return nil
}

// Get retrieves a value from Memcached if it exists and has not expired.
func (m *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := m.client.Get(key)
	if err == memcache.ErrCacheMiss {
		logger.Logger.Warnw("Cache miss", "key", key)
		return nil, nil
	}
	if err != nil {
		logger.Logger.Errorw("Failed to get value from Memcached", "key", key, "error", err)
		return nil, err
	}

	return string(item.Value), nil
}
