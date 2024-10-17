package configs

// GetAppConfig returns application configurations as a map, similar to Laravel's configuration files.
func GetAppConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":            GetWithDefault("APP_NAME", "JazzApp"),
		"env":             GetWithDefault("APP_ENV", "production"),
		"debug":           GetWithDefault("APP_DEBUG", false),
		"url":             GetWithDefault("APP_URL", "http://localhost"),
		"timezone":        GetWithDefault("APP_TIMEZONE", "UTC"),
		"locale":          GetWithDefault("APP_LOCALE", "en"),
		"fallback_locale": GetWithDefault("APP_FALLBACK_LOCALE", "en"),
		"faker_locale":    GetWithDefault("APP_FAKER_LOCALE", "en_US"),
		"cipher":          GetWithDefault("APP_CIPHER", "AES-256-CBC"),
		"key":             Get("APP_KEY"),
		"previous_keys":   GetWithDefault("APP_PREVIOUS_KEYS", ","),
		"maintenance": map[string]interface{}{
			"driver": GetWithDefault("APP_MAINTENANCE_DRIVER", "file"),
			"store":  GetWithDefault("APP_MAINTENANCE_STORE", "database"),
		},
	}
}
