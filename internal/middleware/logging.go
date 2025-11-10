package middleware

import (
	"context"
	"time"

	"gin-crud-api/internal/logger"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// RequestIDKey is the context key for request ID
type contextKey string

const RequestIDKey contextKey = "request_id"

// LoggingMiddleware creates a GraphQL operation middleware for structured logging
// Logs each GraphQL operation (query/mutation) with timing and error information
func LoggingMiddleware() graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		// Generate request ID for tracing
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		// Get operation context
		oc := graphql.GetOperationContext(ctx)

		// Create logger with request ID
		log := logger.WithRequestID(requestID)

		// Log operation start
		log.Info().
			Str("operation", oc.OperationName).
			Str("query", oc.RawQuery).
			Msg("GraphQL operation started")

		// Record start time
		start := time.Now()

		// Execute the operation
		response := next(ctx)

		// Calculate duration
		duration := time.Since(start)

		// Log operation completion
		return func(ctx context.Context) *graphql.Response {
			res := response(ctx)

			// Determine log level based on errors
			var logEvent *zerolog.Event
			if res.Errors != nil && len(res.Errors) > 0 {
				logEvent = log.Error()
				for _, err := range res.Errors {
					logEvent = logEvent.
						Str("error", err.Message).
						Str("path", err.Path.String())
				}
			} else {
				logEvent = log.Info()
			}

			// Log operation result
			logEvent.
				Str("operation", oc.OperationName).
				Dur("duration_ms", duration).
				Int("error_count", len(res.Errors)).
				Msg("GraphQL operation completed")

			return res
		}
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
