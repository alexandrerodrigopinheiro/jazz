// logger/logger.go
package logger

import (
	"jazz/backend/configs"

	"go.uber.org/zap"
)

// AppLogger represents the logger instance.
type AppLogger struct {
	sugarLogger *zap.SugaredLogger
}

// Logger is a global logger instance that can be used throughout the application.
var Logger *AppLogger

// InitializeLogger initializes the logger with different log levels for development and production.
func InitializeLogger() {
	var zapLogger *zap.Logger
	var err error

	// Get application configuration using configs module
	appConfig := configs.GetAppConfig()
	env, ok := appConfig["env"].(string)
	if !ok || env == "" {
		env = "development" // Default to "development" if the environment is not set
	}

	if env == "production" {
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout", "../storage/logs/app.log"}
		zapLogger, err = config.Build()
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.LevelKey = "level"
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.CallerKey = "caller"
		config.OutputPaths = []string{"stdout", "../storage/logs/app.log"}
		zapLogger, err = config.Build()
	}

	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	Logger = &AppLogger{
		sugarLogger: zapLogger.Sugar(),
	}
}

// Info logs an info-level message.
func (l *AppLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof logs a formatted info-level message.
func (l *AppLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn logs a warning-level message.
func (l *AppLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// Warnf logs a formatted warning-level message.
func (l *AppLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Errorf logs a formatted error-level message.
func (l *AppLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// Fatal logs a fatal-level message.
func (l *AppLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}
