package cache

import (
	"sync"
	"time"
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
	var result interface{}

	// Check if value is already in the cache
	entry, found := o.cache.Load(key)
	if found {
		cacheEntry := entry.(SwingCacheEntry[interface{}])
		if time.Now().Unix() < cacheEntry.Expiration {
			return cacheEntry.Value, nil
		}
		o.Forget(key)
	}

	// If value is not cached, execute the callback
	result, err := callback()
	if err != nil {
		return result, err
	}

	// Cache the value
	if err := o.Set(key, result, expiration); err != nil {
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
		return nil, nil
	}
	cacheEntry := entry.(SwingCacheEntry[interface{}])
	if time.Now().Unix() > cacheEntry.Expiration {
		o.Forget(key)
		return nil, nil
	}
	return cacheEntry.Value, nil
}
