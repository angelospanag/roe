package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const loggerKey contextKey = "logger"

// RequestID generates a UUIDv7 per request, attaches a logger with the
// request_id field to the context, and passes it downstream.
func RequestID(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, _ := uuid.NewV7()
			ctx := context.WithValue(r.Context(), loggerKey, base.With("request_id", id.String()))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LoggerFromContext returns the request-scoped logger, falling back to the default.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
