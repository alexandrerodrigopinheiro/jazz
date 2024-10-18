package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"jazz/backend/pkg/database"
	"jazz/backend/pkg/logger"

	"gorm.io/gorm"
)

// CacheEntry represents a cache entry in the database.
type CacheEntry struct {
	Key        string `gorm:"primaryKey"`
	Value      string
	Expiration int64
}

// DatabaseCache is a cache implementation that stores values in a database.
type DatabaseCache struct {
	db *gorm.DB
}

// NewDatabaseCache creates a new instance of Cache (DatabaseCache).
func NewDatabaseCache() Cache {
	// Initialize the logger first
	logger.InitializeLogger()

	// Get the database instance from the database module
	db := database.GetDBInstance()

	// Ensure db is not nil
	if db == nil {
		logger.Logger.Fatal("Database connection is nil after initialization.")
	}

	// Check if the table exists and create it if necessary
	if !db.Migrator().HasTable(&CacheEntry{}) {
		logger.Logger.Info("Table 'cache_entries' does not exist. Creating it now...")
		err := db.Migrator().CreateTable(&CacheEntry{})
		if err != nil {
			logger.Logger.Fatal(fmt.Sprintf("Failed to create table 'cache_entries': %v", err))
		}
	}

	// Perform migration to ensure table structure is updated if needed
	err := db.AutoMigrate(&CacheEntry{})
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to migrate table 'cache_entries': %v", err))
	}

	logger.Logger.Info("Database connection successfully established for cache")
	return &DatabaseCache{db: db}
}

// Set stores a value in the database.
func (d *DatabaseCache) Set(key string, value interface{}, expiration time.Duration) error {
	if d.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	valueBytes, err := json.Marshal(value)
	if err != nil {
		logger.Logger.Errorw("Failed to serialize value", "key", key, "error", err)
		return err
	}

	expirationTime := time.Now().Add(expiration).Unix()
	entry := CacheEntry{Key: key, Value: string(valueBytes), Expiration: expirationTime}
	if err := d.db.Save(&entry).Error; err != nil {
		logger.Logger.Errorw("Failed to save value in database", "key", key, "error", err)
		return err
	}

	return nil
}

// Forget removes a value from the cache.
func (d *DatabaseCache) Forget(key string) error {
	if d.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := d.db.Delete(&CacheEntry{}, "`key` = ?", key).Error; err != nil { // Escapando `key` com crases
		logger.Logger.Errorw("Failed to delete value from database cache", "key", key, "error", err)
		return err
	}

	return nil
}

// Remember retrieves a value or executes a callback if not present.
func (d *DatabaseCache) Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	value, err := d.Get(key)
	if err == nil && value != nil {
		logger.Logger.Infow("Cache hit", "key", key)
		return value, nil
	}

	result, err := callback()
	if err != nil {
		logger.Logger.Errorw("Callback execution failed", "key", key, "error", err)
		return nil, err
	}

	if err := d.Set(key, result, expiration); err != nil {
		logger.Logger.Errorw("Failed to set value in cache after callback", "key", key, "error", err)
		return result, err
	}

	return result, nil
}

// Get retrieves a value from the cache if it exists and has not expired.
func (d *DatabaseCache) Get(key string) (interface{}, error) {
	if d.db == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	var entry CacheEntry
	result := d.db.First(&entry, "`key` = ?", key)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Logger.Warnw("Cache miss", "key", key)
			return nil, nil
		}
		logger.Logger.Errorw("Failed to get value from database", "key", key, "error", result.Error)
		return nil, result.Error
	}

	if time.Now().Unix() > entry.Expiration {
		d.Forget(key)
		logger.Logger.Warnw("Cache entry expired", "key", key)
		return nil, nil
	}

	var value interface{}
	if err := json.Unmarshal([]byte(entry.Value), &value); err != nil {
		logger.Logger.Errorw("Failed to deserialize value", "key", key, "error", err)
		return nil, err
	}

	return value, nil
}
