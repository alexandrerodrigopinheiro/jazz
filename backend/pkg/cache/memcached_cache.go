package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"jazz/backend/configs"

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
		fmt.Println("MEMCACHED_HOST or MEMCACHED_PORT is not set or invalid. Falling back to default cache.")
		return nil
	}

	portStr := fmt.Sprintf("%d", port)

	client := memcache.New(fmt.Sprintf("%s:%s", host, portStr))

	// Testing the connection with Memcached
	// `memcache.Client` does not have a `Ping()` function, so we will do a simple `Set` and `Get` operation to test
	testKey := "test_connection"
	testValue := []byte("ping")
	err := client.Set(&memcache.Item{Key: testKey, Value: testValue, Expiration: 1})
	if err != nil {
		fmt.Printf("Memcached unavailable at %s:%s. Error: %s. Falling back to default cache.\n", host, portStr, err)
		return nil
	}

	_, err = client.Get(testKey)
	if err != nil {
		fmt.Printf("Memcached unavailable at %s:%s. Error: %s. Falling back to default cache.\n", host, portStr, err)
		return nil
	}

	fmt.Printf("Connected to Memcached at %s:%s\n", host, portStr)
	return &MemcachedCache{client: client}
}

// Set stores a value in Memcached.
func (m *MemcachedCache) Set(key string, value interface{}, expiration time.Duration) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
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
	var result interface{}

	// Check if value is already in the cache
	value, err := m.Get(key)
	if err == nil && value != nil {
		err = json.Unmarshal([]byte(value.(string)), &result)
		if err == nil {
			return result, nil
		}
	}

	// If value is not cached, execute the callback
	result, err = callback()
	if err != nil {
		return result, err
	}

	// Cache the value
	if err := m.Set(key, result, expiration); err != nil {
		return result, err
	}

	return result, nil
}

// Forget removes a value from Memcached.
func (m *MemcachedCache) Forget(key string) error {
	return m.client.Delete(key)
}

// Get retrieves a value from Memcached if it exists and has not expired.
func (m *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := m.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return string(item.Value), nil
}
