package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

const (
	// loggerKey points to the value in the context where the logger is stored.
	loggerKey = contextKey("logger")
	// defaultLogLevel sets the default logger verbosity level.
	defaultLogLevel = slog.LevelInfo
	// defaultLogJSON sets the default logger output format.
	defaultLogJSON = true
)

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include when calling [DefaultLogger].
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

// NewLogger creates a [slog.Logger] with the given configuration.
func NewLogger(w io.Writer, level slog.Level, json bool) *slog.Logger {
	var handler slog.Handler
	options := &slog.HandlerOptions{Level: level}

	if json {
		handler = slog.NewJSONHandler(w, options)
	} else {
		handler = slog.NewTextHandler(w, options)
	}

	return slog.New(handler)
}

// NewLoggerTimeless creates a [slog.Logger] with the given configuration. It
// does not include the top level time attribute. This is useful for when the
// output must be deterministic, like testing.
func NewLoggerTimeless(w io.Writer, level slog.Level, json bool) *slog.Logger {
	var handler slog.Handler
	options := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	}

	if json {
		handler = slog.NewJSONHandler(w, options)
	} else {
		handler = slog.NewTextHandler(w, options)
	}

	return slog.New(handler)
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger(os.Stderr, defaultLogLevel, defaultLogJSON)
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in a context. If no such context
// exists, a default logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}

// SlogLevel converts the given string to the appropriate log level. The
// supported input options are "info", "warn", "error", and "debug". All
// other inputs will result in an info level. The input is case insensitive.
func SlogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
