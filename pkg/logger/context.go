package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

func WithContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) zerolog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(zerolog.Logger); ok {
		return logger
	}
	return Logger
}

func WithFields(fields map[string]any) zerolog.Logger {
	ctx := Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return ctx.Logger()
}
