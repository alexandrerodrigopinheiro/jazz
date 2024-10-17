// Package Cache - Clean Architecture and Best Practices

package cache

import (
	"fmt"

	"jazz/backend/configs"
	"jazz/backend/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
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

// NewDatabaseCache creates a new instance of DatabaseCache.
func NewDatabaseCache() *DatabaseCache {
	// Get database configuration
	dbConfig := configs.GetDatabaseConfig()
	defaultConnName := dbConfig["default"].(string)
	connections := dbConfig["connections"].(map[string]interface{})
	connectionConfig, ok := connections[defaultConnName].(map[string]interface{})
	if !ok {
		logger.Logger.Fatal(fmt.Sprintf("No configuration found for the database connection: %s", defaultConnName))
	}

	driver := connectionConfig["driver"].(string)
	// Using custom GormLogger instead of Gorm's default logger
	dbLogger := logger.NewGormLogger(logger.Logger, gormLogger.Info)
	var db *gorm.DB
	var err error

	// Initialize the database based on the driver
	switch driver {
	case "sqlite":
		dbPath := "../database/database.sqlite"
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: dbLogger})

	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			connectionConfig["username"].(string),
			connectionConfig["password"].(string),
			connectionConfig["host"].(string),
			connectionConfig["port"].(string),
			connectionConfig["database"].(string),
			connectionConfig["charset"].(string),
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})

	case "pgsql":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			connectionConfig["host"].(string),
			connectionConfig["port"].(string),
			connectionConfig["username"].(string),
			connectionConfig["database"].(string),
			connectionConfig["password"].(string),
			connectionConfig["sslmode"].(string),
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger})

	case "sqlsrv":
		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s",
			connectionConfig["username"].(string),
			connectionConfig["password"].(string),
			connectionConfig["host"].(string),
			connectionConfig["port"].(string),
			connectionConfig["database"].(string),
		)
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{Logger: dbLogger})

	default:
		logger.Logger.Fatal(fmt.Sprintf("Unsupported database driver: %s", driver))
	}

	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Migrate the table only if it does not exist
	if !db.Migrator().HasTable(&CacheEntry{}) {
		err = db.AutoMigrate(&CacheEntry{})
		if err != nil {
			logger.Logger.Fatal(fmt.Sprintf("Failed to migrate database: %v", err))
		}
	}

	logger.Logger.Info("Database connection successfully established")
	return &DatabaseCache{db: db}
}
