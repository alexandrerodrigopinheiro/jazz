// Package Cache - File-based Cache Implementation with Logging
package cache

import (
	"encoding/json"
	"fmt"
	"jazz/backend/pkg/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileCache is a simple file-based cache that stores values in files.
type FileCache struct {
	cacheDir string
}

// NewFileCache creates a new instance of FileCache.
func NewFileCache() *FileCache {
	cacheDir := "../../storage/cache"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			logger.Logger.Errorw("Failed to create cache directory", "error", err)
			return nil
		}
	}
	return &FileCache{cacheDir: cacheDir}
}

// Set stores a value in a file.
func (t *FileCache) Set(key string, value interface{}, expiration time.Duration) error {
	expirationTime := time.Now().Add(expiration).Unix()
	valueBytes, err := json.Marshal(value)
	if err != nil {
		logger.Logger.Errorw("Failed to serialize value", "key", key, "error", err)
		return err
	}

	fileContent := fmt.Sprintf("%d\n%s", expirationTime, string(valueBytes))
	filePath := filepath.Join(t.cacheDir, key+".txt")

	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		logger.Logger.Errorw("Failed to write cache file", "key", key, "error", err)
		return err
	}

	return nil
}

// Remember retrieves a value from the cache or executes a callback to get it if not present.
func (t *FileCache) Remember(key string, expiration time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	var result interface{}

	// Check if value is already in the cache
	value, err := t.Get(key)
	if err == nil && value != "" {
		lines := strings.Split(value.(string), "\n")
		if len(lines) > 1 {
			expirationTime, err := strconv.ParseInt(lines[0], 10, 64)
			if err == nil && time.Now().Unix() < expirationTime {
				err = json.Unmarshal([]byte(lines[1]), &result)
				if err == nil {
					logger.Logger.Infow("Cache hit", "key", key)
					return result, nil
				}
			}
		}
	}

	// If value is not cached, execute the callback
	result, err = callback()
	if err != nil {
		logger.Logger.Errorw("Callback execution failed", "key", key, "error", err)
		return result, err
	}

	// Cache the value
	if err := t.Set(key, result, expiration); err != nil {
		logger.Logger.Errorw("Failed to set value in cache after callback", "key", key, "error", err)
		return result, err
	}

	return result, nil
}

// Forget removes a value from the cache.
func (t *FileCache) Forget(key string) error {
	filePath := filepath.Join(t.cacheDir, key+".txt")
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			logger.Logger.Errorw("Failed to remove cache file", "key", key, "error", err)
			return err
		}
	}
	return nil
}

// Get retrieves a value from the cache if it exists and has not expired.
func (t *FileCache) Get(key string) (interface{}, error) {
	filePath := filepath.Join(t.cacheDir, key+".txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Logger.Warnw("Cache file does not exist", "key", key)
			return nil, nil
		}
		logger.Logger.Errorw("Failed to read cache file", "key", key, "error", err)
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 1 {
		expirationTime, err := strconv.ParseInt(lines[0], 10, 64)
		if err == nil && time.Now().Unix() < expirationTime {
			var value interface{}
			if err := json.Unmarshal([]byte(lines[1]), &value); err == nil {
				logger.Logger.Infow("Cache hit", "key", key)
				return value, nil
			}
		}
		t.Forget(key)
		logger.Logger.Warnw("Cache entry expired", "key", key)
	}

	return nil, nil
}
