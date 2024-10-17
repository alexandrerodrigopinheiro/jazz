package configs

import "strings"

// GetCacheConfig returns cache configurations as a map, similar to Laravel's cache configuration files.
func GetCacheConfig() map[string]interface{} {
	return map[string]interface{}{
		"default": GetWithDefault("CACHE_STORE", "database"),
		"stores": map[string]interface{}{
			"array": map[string]interface{}{
				"driver":    "array",
				"serialize": false,
			},
			"database": map[string]interface{}{
				"driver":          "database",
				"connection":      Get("DB_CACHE_CONNECTION"),
				"table":           GetWithDefault("DB_CACHE_TABLE", "cache"),
				"lock_connection": Get("DB_CACHE_LOCK_CONNECTION"),
				"lock_table":      Get("DB_CACHE_LOCK_TABLE"),
			},
			"file": map[string]interface{}{
				"driver":    "file",
				"path":      GetWithDefault("CACHE_FILE_PATH", "storage/framework/cache/data"),
				"lock_path": GetWithDefault("CACHE_LOCK_PATH", "storage/framework/cache/data"),
			},
			"memcached": map[string]interface{}{
				"driver":        "memcached",
				"persistent_id": Get("MEMCACHED_PERSISTENT_ID"),
				"sasl": []string{
					Get("MEMCACHED_USERNAME").(string),
					Get("MEMCACHED_PASSWORD").(string),
				},
				"options": map[string]interface{}{},
				"servers": []map[string]interface{}{
					{
						"host":   GetWithDefault("MEMCACHED_HOST", "127.0.0.1"),
						"port":   GetWithDefault("MEMCACHED_PORT", 11211),
						"weight": 100,
					},
				},
			},
			"redis": map[string]interface{}{
				"driver":          "redis",
				"connection":      GetWithDefault("REDIS_CACHE_CONNECTION", "cache"),
				"lock_connection": GetWithDefault("REDIS_CACHE_LOCK_CONNECTION", "default"),
				"url":             GetWithDefault("REDIS_URL", "redis://localhost:6379"),
			},
			"dynamodb": map[string]interface{}{
				"driver":   "dynamodb",
				"key":      Get("AWS_ACCESS_KEY_ID"),
				"secret":   Get("AWS_SECRET_ACCESS_KEY"),
				"region":   GetWithDefault("AWS_DEFAULT_REGION", "us-east-1"),
				"table":    GetWithDefault("DYNAMODB_CACHE_TABLE", "cache"),
				"endpoint": Get("DYNAMODB_ENDPOINT"),
			},
			"swing": map[string]interface{}{
				"driver": "swing",
			},
		},
		"prefix": strings.ToLower(strings.ReplaceAll(GetWithDefault("CACHE_PREFIX", Get("APP_NAME").(string)+"_cache_").(string), " ", "_")),
	}
}
