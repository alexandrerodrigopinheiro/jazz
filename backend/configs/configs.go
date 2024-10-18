// configs/configs.go
package configs

import (
	"fmt"
	"jazz/backend/pkg/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	configValues map[string]interface{}
	once         sync.Once
)

// LoadConfig loads environment variables from the .env file using a Singleton pattern.
func LoadConfig() map[string]interface{} {
	once.Do(func() {
		// Initialize the logger first
		logger.InitializeLogger()

		// Get the current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			logger.Logger.Fatalw("Failed to get current working directory", "error", err)
			return
		}

		// Construct the absolute path to the .env file
		envPath := filepath.Join(currentDir, "../../../.env")

		// Load the .env file
		err = godotenv.Load(envPath)
		if err != nil {
			logger.Logger.Fatalw("Failed to load .env file", "path", envPath, "error", err)
			return
		}

		// Confirm that the environment variables have been loaded by checking a known variable
		requiredEnvVar := "APP_ENV"
		if os.Getenv(requiredEnvVar) == "" {
			logger.Logger.Warnw("Environment variable not loaded from .env file", "key", requiredEnvVar)
			logger.Logger.Fatal(fmt.Errorf("environment variable %s not loaded from .env file", requiredEnvVar))
			return
		}

		// Load environment variables into the configValues map
		configValues = make(map[string]interface{})
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			if len(pair) == 2 {
				configValues[pair[0]] = pair[1]
			}
		}
	})

	return configValues
}

// Get returns the value of a specific environment variable.
func Get(key string) interface{} {
	configs := LoadConfig()
	if value, exists := configs[key]; exists {
		return value
	}
	return ""
}

// GetWithDefault returns the value of a specific environment variable or a default value if not set.
func GetWithDefault(key string, defaultValue interface{}) interface{} {
	configs := LoadConfig()
	if value, exists := configs[key]; exists {
		return value
	}
	return defaultValue
}

// All returns a map with all loaded environment variables.
func All() map[string]interface{} {
	return LoadConfig()
}
