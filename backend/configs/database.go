// Package configs - Database configurations for the Jazz framework

package configs

import (
	"strings"
)

// GetDatabaseConfig returns database configurations as a map, similar to Laravel's database configuration file.
func GetDatabaseConfig() map[string]interface{} {
	return map[string]interface{}{
		"default": GetWithDefault("DB_CONNECTION", "sqlite"),
		"connections": map[string]interface{}{
			"sqlite": map[string]interface{}{
				"driver":                  "sqlite",
				"url":                     Get("DB_URL"),
				"database":                GetWithDefault("DB_DATABASE", "database/database.sqlite"),
				"prefix":                  "",
				"foreign_key_constraints": GetWithDefault("DB_FOREIGN_KEYS", true),
				"busy_timeout":            Get("DB_BUSY_TIMEOUT"),
				"journal_mode":            Get("DB_JOURNAL_MODE"),
				"synchronous":             Get("DB_SYNCHRONOUS"),
			},
			"mysql": map[string]interface{}{
				"driver":         "mysql",
				"url":            Get("DB_URL"),
				"host":           GetWithDefault("DB_HOST", "127.0.0.1"),
				"port":           GetWithDefault("DB_PORT", "3306"),
				"database":       GetWithDefault("DB_DATABASE", "jazz"),
				"username":       GetWithDefault("DB_USERNAME", "root"),
				"password":       Get("DB_PASSWORD"),
				"unix_socket":    Get("DB_SOCKET"),
				"charset":        GetWithDefault("DB_CHARSET", "utf8mb4"),
				"collation":      GetWithDefault("DB_COLLATION", "utf8mb4_unicode_ci"),
				"prefix":         "",
				"prefix_indexes": true,
				"strict":         true,
				"engine":         Get("DB_ENGINE"),
				"options":        map[string]interface{}{"ssl_ca": Get("MYSQL_ATTR_SSL_CA")},
			},
			"mariadb": map[string]interface{}{
				"driver":         "mariadb",
				"url":            Get("DB_URL"),
				"host":           GetWithDefault("DB_HOST", "127.0.0.1"),
				"port":           GetWithDefault("DB_PORT", "3306"),
				"database":       GetWithDefault("DB_DATABASE", "jazz"),
				"username":       GetWithDefault("DB_USERNAME", "root"),
				"password":       Get("DB_PASSWORD"),
				"unix_socket":    Get("DB_SOCKET"),
				"charset":        GetWithDefault("DB_CHARSET", "utf8mb4"),
				"collation":      GetWithDefault("DB_COLLATION", "utf8mb4_unicode_ci"),
				"prefix":         "",
				"prefix_indexes": true,
				"strict":         true,
				"engine":         Get("DB_ENGINE"),
				"options":        map[string]interface{}{"ssl_ca": Get("MYSQL_ATTR_SSL_CA")},
			},
			"pgsql": map[string]interface{}{
				"driver":         "pgsql",
				"url":            Get("DB_URL"),
				"host":           GetWithDefault("DB_HOST", "127.0.0.1"),
				"port":           GetWithDefault("DB_PORT", "5432"),
				"database":       GetWithDefault("DB_DATABASE", "jazz"),
				"username":       GetWithDefault("DB_USERNAME", "root"),
				"password":       Get("DB_PASSWORD"),
				"charset":        GetWithDefault("DB_CHARSET", "utf8"),
				"prefix":         "",
				"prefix_indexes": true,
				"search_path":    GetWithDefault("DB_SEARCH_PATH", "public"),
				"sslmode":        GetWithDefault("DB_SSLMODE", "prefer"),
			},
			"sqlsrv": map[string]interface{}{
				"driver":                   "sqlsrv",
				"url":                      Get("DB_URL"),
				"host":                     GetWithDefault("DB_HOST", "localhost"),
				"port":                     GetWithDefault("DB_PORT", "1433"),
				"database":                 GetWithDefault("DB_DATABASE", "jazz"),
				"username":                 GetWithDefault("DB_USERNAME", "root"),
				"password":                 Get("DB_PASSWORD"),
				"charset":                  GetWithDefault("DB_CHARSET", "utf8"),
				"prefix":                   "",
				"prefix_indexes":           true,
				"encrypt":                  Get("DB_ENCRYPT"),
				"trust_server_certificate": Get("DB_TRUST_SERVER_CERTIFICATE"),
			},
		},
		"migrations": map[string]interface{}{
			"table":                  GetWithDefault("DB_MIGRATIONS_TABLE", "migrations"),
			"update_date_on_publish": GetWithDefault("DB_UPDATE_DATE_ON_PUBLISH", true),
		},
		"redis": map[string]interface{}{
			"client": GetWithDefault("REDIS_CLIENT", "redis"),
			"options": map[string]interface{}{
				"cluster": GetWithDefault("REDIS_CLUSTER", "redis"),
				"prefix":  strings.ToLower(strings.ReplaceAll(GetWithDefault("CACHE_PREFIX", GetWithDefault("APP_NAME", "jazz").(string)+"_cache_").(string), " ", "_")),
			},
			"default": map[string]interface{}{
				"url":      Get("REDIS_URL"),
				"host":     GetWithDefault("REDIS_HOST", "127.0.0.1"),
				"username": Get("REDIS_USERNAME"),
				"password": Get("REDIS_PASSWORD"),
				"port":     GetWithDefault("REDIS_PORT", "6379"),
				"database": GetWithDefault("REDIS_DB", "0"),
			},
			"cache": map[string]interface{}{
				"url":      Get("REDIS_URL"),
				"host":     GetWithDefault("REDIS_HOST", "127.0.0.1"),
				"username": Get("REDIS_USERNAME"),
				"password": Get("REDIS_PASSWORD"),
				"port":     GetWithDefault("REDIS_PORT", "6379"),
				"database": GetWithDefault("REDIS_CACHE_DB", "1"),
			},
		},
	}
}
