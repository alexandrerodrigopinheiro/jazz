package cache

import (
	"jazz/backend/pkg/logger"
	"testing"
	"time"
)

// Helper function to run common cache tests on different cache implementations
func runCommonCacheTests(t *testing.T, cache Cache) {
	key := "test_key"
	value := "test_value"
	expiration := time.Second * 10

	// Test Set
	err := cache.Set(key, value, expiration)
	if err != nil {
		logger.Logger.Errorf("Failed to set value in cache: %s", err)
		t.Fatalf("Failed to set value in cache: %s", err)
	}

	// Test Get
	cachedValue, err := cache.Get(key)
	if err != nil {
		logger.Logger.Errorf("Failed to get value from cache: %s", err)
		t.Fatalf("Failed to get value from cache: %s", err)
	}
	if cachedValue != value {
		logger.Logger.Warnf("Expected %s but got %s", value, cachedValue)
		t.Errorf("Expected %s but got %s", value, cachedValue)
	}

	// Test Forget
	err = cache.Forget(key)
	if err != nil {
		logger.Logger.Errorf("Failed to forget value in cache: %s", err)
		t.Fatalf("Failed to forget value in cache: %s", err)
	}

	// Ensure value is removed
	cachedValue, err = cache.Get(key)
	if err != nil {
		logger.Logger.Errorf("Failed to get value from cache: %s", err)
		t.Fatalf("Failed to get value from cache: %s", err)
	}
	if cachedValue != nil {
		logger.Logger.Warn("Expected value to be removed from cache")
		t.Errorf("Expected value to be removed from cache")
	}
}

// TestFileCache tests the FileCache implementation
func TestFileCache(t *testing.T) {
	cache := NewFileCache()
	if cache == nil {
		logger.Logger.Warn("Skipping FileCache: FileCache is not available")
		t.Skip("Skipping FileCache: FileCache is not available")
	}
	runCommonCacheTests(t, cache)
}

// TestRedisCache tests the RedisCache implementation
func TestRedisCache(t *testing.T) {
	cache := NewRedisCache()
	if cache == nil {
		logger.Logger.Warn("Skipping RedisCache: Redis is not available or not configured properly")
		t.Skip("Skipping RedisCache: Redis is not available or not configured properly")
	}
	runCommonCacheTests(t, cache)
}

// TestDatabaseCache tests the DatabaseCache implementation
func TestDatabaseCache(t *testing.T) {
	cache := NewDatabaseCache()
	if cache == nil {
		logger.Logger.Warn("Skipping DatabaseCache: Database is not available or not configured properly")
		t.Skip("Skipping DatabaseCache: Database is not available or not configured properly")
	}
	runCommonCacheTests(t, cache)
}

// TestDynamoDBCache tests the DynamoDBCache implementation
func TestDynamoDBCache(t *testing.T) {
	cache := NewDynamoDBCache()
	if cache == nil {
		logger.Logger.Warn("Skipping DynamoDBCache: DynamoDB is not available or credentials are missing")
		t.Skip("Skipping DynamoDBCache: DynamoDB is not available or credentials are missing")
	}
	runCommonCacheTests(t, cache)
}

// TestMemcachedCache tests the MemcachedCache implementation
func TestMemcachedCache(t *testing.T) {
	cache := NewMemcachedCache()
	if cache == nil {
		logger.Logger.Warn("Skipping MemcachedCache: Memcached is not available or not configured properly")
		t.Skip("Skipping MemcachedCache: Memcached is not available or not configured properly")
	}
	runCommonCacheTests(t, cache)
}

// TestSwingCache tests the SwingCache implementation
func TestSwingCache(t *testing.T) {
	cache := NewSwingCache()
	if cache == nil {
		logger.Logger.Warn("Skipping SwingCache: SwingCache is not available")
		t.Skip("Skipping SwingCache: SwingCache is not available")
	}
	runCommonCacheTests(t, cache)
}
