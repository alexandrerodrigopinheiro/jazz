package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// AppLogger represents the logger instance.
type AppLogger struct {
	sugarLogger *zap.SugaredLogger
}

// Logger is a global logger instance that can be used throughout the application.
var Logger *AppLogger
var isLoggerInitialized bool

// InitializeLogger initializes the logger with different log levels for development and production.
func InitializeLogger() {
	if isLoggerInitialized {
		return
	}

	// Ensure the logs directory exists
	logsPath := "../../storage/logs"
	if err := os.MkdirAll(logsPath, os.ModePerm); err != nil {
		panic("failed to create logs directory: " + err.Error())
	}

	var zapLogger *zap.Logger
	var err error

	// Get application environment variable
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	logFilePath := filepath.Join(logsPath, "app.log")

	if env == "production" {
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout", logFilePath}
		zapLogger, err = config.Build()
	} else {
		config := zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout", logFilePath}
		zapLogger, err = config.Build()
	}

	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	Logger = &AppLogger{
		sugarLogger: zapLogger.Sugar(),
	}

	isLoggerInitialized = true
}

// Resto do código...

// GetLogger retorna a instância singleton do logger, inicializando-a se necessário.
func GetLogger() *AppLogger {
	InitializeLogger()
	return Logger
}

// Métodos para registrar logs em diferentes níveis

// Info logs an info-level message.
func (l *AppLogger) Info(args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Info(args...)
}

// Infof logs a formatted info-level message.
func (l *AppLogger) Infof(template string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Infof(template, args...)
}

// Infow logs a message info-level message.
func (l *AppLogger) Infow(msg string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Infow(msg, args...)
}

// Warn logs a warning-level message.
func (l *AppLogger) Warn(args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Warn(args...)
}

// Warnf logs a formatted warning-level message.
func (l *AppLogger) Warnf(template string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Warnf(template, args...)
}

// Warnw logs a message warning-level message.
func (l *AppLogger) Warnw(msg string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Warnw(msg, args...)
}

// Error logs an error-level message.
func (l *AppLogger) Error(args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Error(args...)
}

// Errorf logs a formatted error-level message.
func (l *AppLogger) Errorf(template string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Errorf(template, args...)
}

// Errorw logs a message error-level message.
func (l *AppLogger) Errorw(msg string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Errorw(msg, args...)
}

// Debugw logs a debug-level message.
func (l *AppLogger) Debugw(msg string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Debugw(msg, args...)
}

// Fatal logs a fatal-level message.
func (l *AppLogger) Fatal(args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Fatal(args...)
}

// Fatalw logs a formatted fatal-level message.
func (l *AppLogger) Fatalw(template string, args ...interface{}) {
	ensureLoggerInitialized()
	l.sugarLogger.Fatalw(template, args...)
}

// ensureLoggerInitialized garante que o logger tenha sido inicializado antes de acessá-lo.
func ensureLoggerInitialized() {
	if Logger == nil {
		InitializeLogger()
	}
}
