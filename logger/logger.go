package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger ILogger
)

type loggerKey struct{}

type Logger struct {
	l  *zap.Logger
	wl *zap.Logger
}

func (l *Logger) wrapLogger(ctx context.Context) *zap.Logger {
	v := ctx.Value(loggerKey{})
	if v != nil {
		return v.(*zap.Logger)
	}

	return l.wl
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.wrapLogger(ctx).Debug(msg, fields...)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.wrapLogger(ctx).Debug(msg)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.wrapLogger(ctx).Info(msg, fields...)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.wrapLogger(ctx).Info(msg)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.wrapLogger(ctx).Warn(msg, fields...)
}

func (l *Logger) Wanrf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.wrapLogger(ctx).Warn(msg)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.wrapLogger(ctx).Error(msg, fields...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.wrapLogger(ctx).Error(msg)
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.wrapLogger(ctx).Panic(msg, fields...)
}

func (l *Logger) Panicf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.wrapLogger(ctx).Panic(msg)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.wrapLogger(ctx).Fatal(msg, fields...)
}

func (l *Logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.wrapLogger(ctx).Fatal(msg)
}

func (l *Logger) ContextWithFields(ctx context.Context, fields ...zap.Field) context.Context {
	wl := l.wrapLogger(ctx)
	newWL := wl.With(fields...)
	return context.WithValue(ctx, loggerKey{}, newWL)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Debug(ctx, msg, fields...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.Debugf(ctx, format, args)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(ctx, msg, fields...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.Infof(ctx, format, args)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Warn(ctx, msg, fields...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.Wanrf(ctx, format, args)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(ctx, msg, fields...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.Errorf(ctx, format, args)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Panic(ctx, msg, fields...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	logger.Panicf(ctx, format, args)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Fatal(ctx, msg, fields...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logger.Fatalf(ctx, format, args)
}

func ContextWithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return logger.ContextWithFields(ctx, fields...)
}

func NewLogger(config *Config) *Logger {
	if config == nil {
		panic("config is nil")
	}

	config.setDefault()

	core := zapcore.NewCore(Encoder(false),
		zapcore.AddSync(&lumberjack.Logger{
			Filename: filepath.Join(config.FilePath, fmt.Sprintf("%s.log", config.Name)),
			MaxSize:  config.MaxSize,
			MaxAge:   config.MaxAge,
		}),
		config.Level,
	)
	if config.Console {
		core = zapcore.NewTee(core, zapcore.NewCore(Encoder(true), zapcore.AddSync(os.Stdout), config.Level))
	}

	l := zap.New(core, zap.WithCaller(true), zap.AddStacktrace(config.StacktraceLevel))
	wl := l.WithOptions(zap.AddCallerSkip(1))
	for _, f := range config.Fields {
		wl = wl.With(f)
	}
	return &Logger{
		l:  l,
		wl: wl,
	}
}

func Init(config *Config) {
	logger = NewLogger(config)
}

func init() {
	l := zap.New(zapcore.NewCore(Encoder(true), os.Stdout, zapcore.InfoLevel))
	logger = &Logger{
		l:  l,
		wl: l.WithOptions(zap.AddCallerSkip(1)),
	}
}
