package logger

import (
	"context"

	"go.uber.org/zap"
)

type ILogger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Wanrf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Errorf(ctx context.Context, format string, args ...interface{})
	Panic(ctx context.Context, msg string, fields ...zap.Field)
	Panicf(ctx context.Context, format string, args ...interface{})
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Fatalf(ctx context.Context, format string, args ...interface{})
	ContextWithFields(ctx context.Context, fields ...zap.Field) context.Context
}
