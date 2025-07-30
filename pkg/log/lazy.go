package log

import (
	"log/slog"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var (
	loggerCache = make(map[string]*slog.Logger)
	loggerMutex sync.RWMutex
)

// GetPackageLogger returns the logger to use for given packageName, following OpenTelemetry recommendations.
// It creates it if none is present. Beware that calling GetPackageLogger before a [go.opentelemetry.io/otel/log/global]
// [go.opentelemetry.io/otel/log].LoggerProvider is configured will result in an [slog.Logger] that is backed by a noop
// OpenTelemetry logger, that is, a noop logger. This limitation is usually circumvented by calling GetPackageLogger
// after the server is started, itself started after global logger configuration.
// This function is safe for concurrent use.
func GetPackageLogger(packageName string) *slog.Logger {
	loggerMutex.RLock()
	if logger, exists := loggerCache[packageName]; exists {
		loggerMutex.RUnlock()
		return logger
	}
	loggerMutex.RUnlock()

	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// If logger has been created by another routine before this one locked, use it
	logger, exists := loggerCache[packageName]
	if exists {
		return logger
	}

	logger = otelslog.NewLogger(packageName)
	loggerCache[packageName] = logger
	return logger
}
