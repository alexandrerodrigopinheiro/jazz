package cache

import (
	"jazz/backend/pkg/logger"
	"os"
)

// NewCacheManager returns an implementation of Cache based on the configuration driver.
func NewCacheManager() Cache {
	driver := os.Getenv("CACHE_DRIVER")

	switch driver {
	case "redis":
		if cache := NewRedisCache(); cache != nil {
			logger.Logger.Infof("Using Redis Cache as driver: %s", driver)
			return cache
		}
		logger.Logger.Warnf("Redis Cache unavailable, falling back to next available cache for driver: %s", driver)
	case "memcached":
		if cache := NewMemcachedCache(); cache != nil {
			logger.Logger.Infof("Using Memcached Cache as driver: %s", driver)
			return cache
		}
		logger.Logger.Warnf("Memcached Cache unavailable, falling back to next available cache for driver: %s", driver)
	case "database":
		if cache := NewDatabaseCache(); cache != nil {
			logger.Logger.Infof("Using Database Cache as driver: %s", driver)
			return cache
		}
		logger.Logger.Warnf("Database Cache unavailable, falling back to next available cache for driver: %s", driver)
	case "dynamodb":
		if cache := NewDynamoDBCache(); cache != nil {
			logger.Logger.Infof("Using DynamoDB Cache as driver: %s", driver)
			return cache
		}
		logger.Logger.Warnf("DynamoDB Cache unavailable, falling back to next available cache for driver: %s", driver)
	case "swing":
		if cache := NewSwingCache(); cache != nil {
			logger.Logger.Infof("Using Swing Cache (in-memory) as driver: %s", driver)
			return cache
		}
		logger.Logger.Warnf("Swing Cache unavailable, falling back to next available cache for driver: %s", driver)
	default:
		logger.Logger.Warnf("Invalid Cache driver specified: %s, falling back to File Cache", driver)
	}

	// Fallback to file cache
	if cache := NewFileCache(); cache != nil {
		logger.Logger.Info("Using File Cache as default")
		return cache
	}

	// Final fallback to in-memory cache (Swing)
	logger.Logger.Warn("File Cache unavailable, falling back to Swing Cache (in-memory) as final fallback")
	return NewSwingCache()
}
