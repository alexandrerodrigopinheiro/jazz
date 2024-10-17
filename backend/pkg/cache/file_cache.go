package cache

import (
	"encoding/json"
	"fmt"
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
	cacheDir := "../storage/cache"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}
	return &FileCache{cacheDir: cacheDir}
}

// Set stores a value in a file.
func (t *FileCache) Set(key string, value interface{}, expiration time.Duration) error {
	expirationTime := time.Now().Add(expiration).Unix()
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	fileContent := fmt.Sprintf("%d\n%s", expirationTime, string(valueBytes))
	filePath := filepath.Join(t.cacheDir, key+".txt")

	return os.WriteFile(filePath, []byte(fileContent), 0644)
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
					return result, nil
				}
			}
		}
	}

	// If value is not cached, execute the callback
	result, err = callback()
	if err != nil {
		return result, err
	}

	// Cache the value
	if err := t.Set(key, result, expiration); err != nil {
		return result, err
	}

	return result, nil
}

// Forget removes a value from the cache.
func (t *FileCache) Forget(key string) error {
	filePath := filepath.Join(t.cacheDir, key+".txt")
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}
	return nil
}

// Get retrieves a value from the cache if it exists and has not expired.
func (t *FileCache) Get(key string) (interface{}, error) {
	filePath := filepath.Join(t.cacheDir, key+".txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 1 {
		expirationTime, err := strconv.ParseInt(lines[0], 10, 64)
		if err == nil && time.Now().Unix() < expirationTime {
			var value interface{}
			if err := json.Unmarshal([]byte(lines[1]), &value); err == nil {
				return value, nil
			}
		}
		t.Forget(key)
	}

	return nil, nil
}
