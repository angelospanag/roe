package middleware

import (
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
)

// init initializes the default logger using a JSON handler and sets it as the global logger for the application.
func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)
}

// LoggingMiddleware is a Huma middleware function for logging request details and execution duration.
// It generates a unique correlation ID for each request and attaches a logger instance to the context.
func LoggingMiddleware(ctx huma.Context, next func(huma.Context)) {
	start := time.Now()

	// Get the real remote IP address from X-Forwarded-For header or direct connection
	// Try to get X-Forwarded-For header first (used by load balancers)
	rawXFF := ctx.Header("X-Forwarded-For")
	remoteAddress := ""

	if rawXFF != "" {
		// X-Forwarded-For can contain multiple IPs, split them
		parts := strings.Split(rawXFF, ",")
		if len(parts) > 0 {
			// Take the first IP (original client IP) and validate it
			ip := net.ParseIP(strings.TrimSpace(parts[0]))
			if ip != nil {
				remoteAddress = ip.String()
			}
		}
	}

	// If we couldn't get a valid IP from X-Forwarded-For, use RemoteAddr
	if remoteAddress == "" {
		// Try to split host and port from RemoteAddr
		host, _, err := net.SplitHostPort(ctx.RemoteAddr())
		if err == nil {
			// Validate the IP address
			ip := net.ParseIP(host)
			if ip != nil {
				remoteAddress = ip.String()
			} else {
				// If not valid IP but we have a host, use it as a fallback
				remoteAddress = host
			}
		} else {
			// If we can't split host:port, use raw RemoteAddr as a last resort
			remoteAddress = ctx.RemoteAddr()
		}
	}

	// Generate a correlation ID for this request
	correlationID := uuid.Must(uuid.NewV4()).String()
	logger := slog.Default().With(
		slog.String("correlationID", correlationID),
	)
	// Wrap the context to add a value.
	// https://huma.rocks/features/middleware/#context-values
	ctx = huma.WithValue(ctx, "logger", logger)

	next(ctx)

	duration := time.Since(start)
	logger.Info("Request completed",
		slog.String("method", ctx.Method()),
		slog.String("path", ctx.URL().Path),
		slog.Int("status", ctx.Status()),
		slog.String("remoteAddress", remoteAddress),
		slog.Int64("durationMs", duration.Milliseconds()),
	)
}
