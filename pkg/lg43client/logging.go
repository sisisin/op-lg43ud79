package lg43client

import (
	"context"
	"log"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
)

type logLevelKey struct{}

func WithLogLevel(ctx context.Context, level LogLevel) context.Context {
	return context.WithValue(ctx, logLevelKey{}, level)
}

func getLogLevel(ctx context.Context) LogLevel {
	v := ctx.Value(logLevelKey{})
	if v == nil {
		return LogLevelInfo
	}
	return v.(LogLevel)
}

func isDebug(ctx context.Context) bool {
	return getLogLevel(ctx) == LogLevelDebug
}

func isInfo(ctx context.Context) bool {
	return getLogLevel(ctx) <= LogLevelInfo
}

func logDebug(ctx context.Context, format string, v ...any) {
	if isDebug(ctx) {
		log.Printf(format, v...)
	}
}
func logInfo(ctx context.Context, format string, v ...any) {
	if isInfo(ctx) {
		log.Printf(format, v...)
	}
}
