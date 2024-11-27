package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestDefauletLogger(t *testing.T) {
	ctx := context.Background()
	Debug(ctx, "test")
	Debugf(ctx, "test %s", "test")
	Info(ctx, "test")
	Infof(ctx, "test %s", "test")
	Warn(ctx, "test")
	Warnf(ctx, "test %s", "test")
	Error(ctx, "test")
	Errorf(ctx, "test %s", "test")
}

func TestDefauletLoggerWithCtx(t *testing.T) {
	ctx := ContextWithFields(context.Background(), zap.String("key", "value"))
	Debug(ctx, "test")
	Debugf(ctx, "test %s", "test")
	Info(ctx, "test")
	Infof(ctx, "test %s", "test")
	Warn(ctx, "test")
	Warnf(ctx, "test %s", "test")
	Error(ctx, "test")
	Errorf(ctx, "test %s", "test")
}

func TestNewLoggerWithCtx(t *testing.T) {
	Init(&Config{
		Console: true,
		Level:   zapcore.DebugLevel,
		Fields:  []zapcore.Field{zap.String("key1", "value1")},
	})
	ctx := context.Background()
	Debug(ctx, "test")
	Debugf(ctx, "test %s", "test")
	Info(ctx, "test")
	Infof(ctx, "test %s", "test")
	ctx = ContextWithFields(context.Background(), zap.String("key2", "value2"))
	Warn(ctx, "test")
	Warnf(ctx, "test %s", "test")
	Error(ctx, "test")
	Errorf(ctx, "test %s", "test")
}

func TestPanic(t *testing.T) {
	Init(&Config{
		Console: true,
		Level:   zapcore.DebugLevel,
		Fields:  []zapcore.Field{zap.String("key1", "value1")},
	})
	ctx := ContextWithFields(context.Background(), zap.String("key2", "value2"))
	Panic(ctx, "test")
}
