package oplg43ud79

import "context"

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
)

type logLevelKey struct{}

func withLogLevel(ctx context.Context, level LogLevel) context.Context {
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
