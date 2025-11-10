package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Global logger instance
var Logger zerolog.Logger

// Init initializes the global logger with configuration
// level: "debug", "info", "warn", "error" - default is "info"
// pretty: if true, uses human-readable console output (good for development)
func Init(level string, pretty bool) {
	var output io.Writer = os.Stdout

	// Pretty logging for development (with colors and human-readable format)
	if pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// Set log level
	logLevel := zerolog.InfoLevel
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// Create logger with timestamp and level
	Logger = zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Caller(). // Add file:line information
		Logger()

	// Set as global logger
	log.Logger = Logger

	Logger.Info().
		Str("level", logLevel.String()).
		Bool("pretty", pretty).
		Msg("Logger initialized")
}

// GetLogger returns the global logger instance
func GetLogger() zerolog.Logger {
	return Logger
}

// WithRequestID creates a child logger with request ID for tracing
func WithRequestID(requestID string) zerolog.Logger {
	return Logger.With().Str("request_id", requestID).Logger()
}

// WithComponent creates a child logger with component name
func WithComponent(component string) zerolog.Logger {
	return Logger.With().Str("component", component).Logger()
}
