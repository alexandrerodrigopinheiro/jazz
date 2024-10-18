// Package logger - Custom Logger for Gorm integration
package logger

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm/logger"
)

var (
	gormLoggerInstance *GormLogger
	gormLoggerOnce     sync.Once
)

// GormLogger is a custom logger that implements the Gorm logger.Interface.
// It uses the AppLogger instance for logging.
type GormLogger struct {
	logger   *AppLogger
	logLevel logger.LogLevel
}

// NewGormLogger initializes and returns a singleton instance of GormLogger using the global AppLogger.
func NewGormLogger(level logger.LogLevel) *GormLogger {
	gormLoggerOnce.Do(func() {
		gormLoggerInstance = &GormLogger{
			logger:   GetLogger(), // Usa a instância do AppLogger
			logLevel: level,
		}
	})
	return gormLoggerInstance
}

// LogMode sets the logging level for the Gorm logger.
func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.logLevel = level
	return g
}

// Info logs informational messages.
func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel >= logger.Info {
		g.logger.Infof(msg, data...)
	}
}

// Warn logs warning messages.
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel >= logger.Warn {
		g.logger.Warnf(msg, data...)
	}
}

// Error logs error messages.
func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel >= logger.Error {
		g.logger.Errorf(msg, data...)
	}
}

// Trace logs SQL queries along with their execution time.
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && g.logLevel >= logger.Error:
		g.logger.Errorf("SQL Error: %s - %v [Rows affected: %d, Elapsed time: %v]", sql, err, rows, elapsed)
	case elapsed > time.Millisecond*200 && g.logLevel >= logger.Warn:
		// Log if the execution time is longer than 200ms
		g.logger.Warnf("Slow SQL (> 200ms): %s [Rows affected: %d, Elapsed time: %v]", sql, rows, elapsed)
	case g.logLevel >= logger.Info:
		g.logger.Infof("SQL Executed: %s [Rows affected: %d, Elapsed time: %v]", sql, rows, elapsed)
	}
}
