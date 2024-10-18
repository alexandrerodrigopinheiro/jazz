// database/database.go
package database

import (
	"fmt"
	"jazz/backend/configs"
	"jazz/backend/pkg/logger"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	dbInstance        *gorm.DB
	once              sync.Once
	initializationErr error
)

// InitializeDatabase initializes the database using environment variables.
func InitializeDatabase() *gorm.DB {
	once.Do(func() {
		if err := configs.LoadConfig(); err != nil {
			initializationErr = fmt.Errorf("error loading config: %v", err)
			return
		}

		dbConfig := configs.GetDatabaseConfig()

		// Get the default connection name
		defaultConnName, ok := dbConfig["default"].(string)
		if !ok || defaultConnName == "" {
			initializationErr = fmt.Errorf("no default database connection specified in configuration")
			return
		}

		// Get the connection configuration based on the default connection name
		connections, ok := dbConfig["connections"].(map[string]interface{})
		if !ok {
			initializationErr = fmt.Errorf("no database connections found in configuration")
			return
		}

		connectionConfig, ok := connections[defaultConnName].(map[string]interface{})
		if !ok {
			initializationErr = fmt.Errorf("no configuration found for the database connection: %s", defaultConnName)
			return
		}

		driver, ok := connectionConfig["driver"].(string)
		if !ok || driver == "" {
			initializationErr = fmt.Errorf("database driver is not specified in the connection configuration")
			return
		}

		var err error
		switch driver {
		case "mysql", "mariadb":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
				connectionConfig["username"].(string),
				connectionConfig["password"].(string),
				connectionConfig["host"].(string),
				connectionConfig["port"].(string),
				connectionConfig["database"].(string),
				connectionConfig["charset"].(string))
			dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: logger.NewGormLogger(gormLogger.Info),
			})
		case "pgsql":
			dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
				connectionConfig["host"].(string),
				connectionConfig["port"].(string),
				connectionConfig["username"].(string),
				connectionConfig["database"].(string),
				connectionConfig["password"].(string),
				connectionConfig["sslmode"].(string))
			dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.NewGormLogger(gormLogger.Info),
			})
		case "sqlite":
			dbPath, ok := connectionConfig["database"].(string)
			if !ok || dbPath == "" {
				initializationErr = fmt.Errorf("sqlite database path is not specified")
				return
			}
			dbInstance, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
				Logger: logger.NewGormLogger(gormLogger.Info),
			})
		default:
			initializationErr = fmt.Errorf("unsupported database driver: %s", driver)
			return
		}

		if err != nil {
			initializationErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		// Verify if the database connection is valid
		sqlDB, err := dbInstance.DB()
		if err != nil {
			initializationErr = fmt.Errorf("failed to get database instance: %w", err)
			return
		}

		if err = sqlDB.Ping(); err != nil {
			initializationErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		logger.Logger.Info("Database connection successfully established")
	})

	if initializationErr != nil {
		logger.Logger.Fatal(initializationErr.Error())
	}

	return dbInstance
}

// GetDBInstance returns the singleton instance of the database.
func GetDBInstance() *gorm.DB {
	if dbInstance == nil {
		return InitializeDatabase()
	}
	return dbInstance
}
