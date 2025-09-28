package oplg43ud79

import (
	"context"
	"log"
)

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
