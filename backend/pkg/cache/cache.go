package cache

import "time"

// Cache interface defines the required methods for a cache implementation.
type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Forget(key string) error
	Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error)
}
