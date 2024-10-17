package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var configValues map[string]interface{}

// LoadConfig loads environment variables from the .env file.
func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	configValues = make(map[string]interface{})
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			configValues[pair[0]] = pair[1]
		}
	}

	return nil
}

// All returns a map with all loaded environment variables.
func All() map[string]interface{} {
	return configValues
}

// Get returns the value of a specific environment variable.
func Get(key string) interface{} {
	if value, exists := configValues[key]; exists {
		return value
	}
	return ""
}

// GetWithDefault returns the value of a specific environment variable or a default value if not set.
func GetWithDefault(key string, defaultValue interface{}) interface{} {
	if value, exists := configValues[key]; exists {
		return value
	}
	return defaultValue
}


// GetCacheDriver returns the cache driver from the environment variables.
func GetCacheDriver() interface{} {
	return GetWithDefault("CACHE_DRIVER", "text")
}
