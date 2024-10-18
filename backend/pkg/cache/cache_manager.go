package cache

import (
	"jazz/backend/pkg/logger"
	"os"
	"sync"
)

var (
	cacheInstance Cache
	cacheOnce     sync.Once
)

// NewCacheManager returns an implementation of Cache based on the configuration driver.
func NewCacheManager() Cache {
	cacheOnce.Do(func() {
		driver := os.Getenv("CACHE_DRIVER")

		switch driver {
		case "redis":
			cache := NewRedisCache()
			if cache != nil {
				logger.Logger.Infof("Using Redis Cache as driver: %s", driver)
				cacheInstance = cache
				return
			}
			logger.Logger.Warnf("Redis Cache unavailable, falling back to next available cache for driver: %s", driver)
		case "memcached":
			cache := NewMemcachedCache()
			if cache != nil {
				logger.Logger.Infof("Using Memcached Cache as driver: %s", driver)
				cacheInstance = cache
				return
			}
			logger.Logger.Warnf("Memcached Cache unavailable, falling back to next available cache for driver: %s", driver)
		case "database":
			cache := NewDatabaseCache()
			if cache != nil {
				logger.Logger.Infof("Using Database Cache as driver: %s", driver)
				cacheInstance = cache
				return
			}
			logger.Logger.Warnf("Database Cache unavailable, falling back to next available cache for driver: %s", driver)
		case "dynamodb":
			cache := NewDynamoDBCache()
			if cache != nil {
				logger.Logger.Infof("Using DynamoDB Cache as driver: %s", driver)
				cacheInstance = cache
				return
			}
			logger.Logger.Warnf("DynamoDB Cache unavailable, falling back to next available cache for driver: %s", driver)
		case "swing":
			cache := NewSwingCache()
			if cache != nil {
				logger.Logger.Infof("Using Swing Cache (in-memory) as driver: %s", driver)
				cacheInstance = cache
				return
			}
			logger.Logger.Warnf("Swing Cache unavailable, falling back to next available cache for driver: %s", driver)
		default:
			logger.Logger.Warnf("Invalid Cache driver specified: %s, falling back to File Cache", driver)
		}

		// Fallback to file cache
		cache := NewFileCache()
		if cache != nil {
			logger.Logger.Info("Using File Cache as default")
			cacheInstance = cache
			return
		}

		// Final fallback to in-memory cache (Swing)
		logger.Logger.Warn("File Cache unavailable, falling back to Swing Cache (in-memory) as final fallback")
		cacheInstance = NewSwingCache()
	})

	return cacheInstance
}
