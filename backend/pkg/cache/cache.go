package cache

import (
	"time"
)

// Cache Interface - defines methods for cache operations
type Cache interface {
	Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Forget(key string) error
	Get(key string) (interface{}, error)
}
