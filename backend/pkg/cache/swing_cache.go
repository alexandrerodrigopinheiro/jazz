package cache

import (
	"sync"
	"time"

	"jazz/backend/pkg/logger"
)

// SwingCache is a simple in-memory cache that uses a sync.Map to store values.
type SwingCache struct {
	cache sync.Map
}

// SwingCacheEntry is a struct that holds a value and its expiration time.
type SwingCacheEntry[T any] struct {
	Value      T
	Expiration int64
}

// NewSwingCache creates a new instance of SwingCache.
func NewSwingCache() *SwingCache {
	// Garantir que o logger esteja inicializado antes de usá-lo
	logger.InitializeLogger()

	logger.Logger.Info("Initializing SwingCache")
	return &SwingCache{}
}

// Set stores a value in the cache with a specified expiration time.
func (o *SwingCache) Set(key string, value interface{}, expiration time.Duration) error {
	expirationTime := time.Now().Add(expiration).Unix()
	entry := SwingCacheEntry[interface{}]{
		Value:      value,
		Expiration: expirationTime,
	}
	o.cache.Store(key, entry)
	return nil
}

// Remember retrieves a value from the cache or executes a callback to get it if not present.
func (o *SwingCache) Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Check if value is already in the cache
	entry, found := o.cache.Load(key)
	if found {
		cacheEntry := entry.(SwingCacheEntry[interface{}])
		if time.Now().Unix() < cacheEntry.Expiration {
			logger.Logger.Infow("Cache hit in SwingCache", "key", key)
			return cacheEntry.Value, nil
		}
		logger.Logger.Warnw("Cache entry expired in SwingCache", "key", key)
		o.Forget(key)
	}

	// If value is not cached, execute the callback
	result, err := callback()
	if err != nil {
		logger.Logger.Errorw("Callback execution failed", "key", key, "error", err)
		return result, err
	}

	// Cache the value
	if err := o.Set(key, result, expiration); err != nil {
		logger.Logger.Errorw("Failed to set value in SwingCache after callback", "key", key, "error", err)
		return result, err
	}

	return result, nil
}

// Forget removes a value from the cache.
func (o *SwingCache) Forget(key string) error {
	o.cache.Delete(key)
	return nil
}

// Get retrieves a value from the cache if it exists and has not expired.
func (o *SwingCache) Get(key string) (interface{}, error) {
	entry, found := o.cache.Load(key)
	if !found {
		logger.Logger.Warnw("Cache miss in SwingCache", "key", key)
		return nil, nil
	}
	cacheEntry := entry.(SwingCacheEntry[interface{}])
	if time.Now().Unix() > cacheEntry.Expiration {
		logger.Logger.Warnw("Cache entry expired in SwingCache", "key", key)
		o.Forget(key)
		return nil, nil
	}
	return cacheEntry.Value, nil
}
