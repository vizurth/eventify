package logger

import (
	"context"
	"go.uber.org/zap"
)

type key string

const (
	KeyForLogger    key = "logger"
	KeyForRequestID key = "requestID"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, *Logger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return ctx, nil, err
	}
	wrapped := &Logger{l: l}
	ctx = context.WithValue(ctx, KeyForLogger, wrapped)
	return ctx, wrapped, nil
}

// GetLoggerFromCtx безопасно достаёт логгер из контекста
func GetLoggerFromCtx(ctx context.Context) *Logger {
	if l, ok := ctx.Value(KeyForLogger).(*Logger); ok {
		return l
	}
	return nil
}

// GetOrCreateLoggerFromCtx возвращает логгер из контекста или создаёт новый
func GetOrCreateLoggerFromCtx(ctx context.Context) *Logger {
	if l := GetLoggerFromCtx(ctx); l != nil {
		return l
	}
	l, _ := zap.NewProduction()
	return &Logger{l: l}
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(KeyForRequestID) != nil {
		fields = append(fields, zap.String(string(KeyForRequestID), ctx.Value(KeyForRequestID).(string)))
	}
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(KeyForRequestID) != nil {
		fields = append(fields, zap.String(string(KeyForRequestID), ctx.Value(KeyForRequestID).(string)))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(KeyForRequestID) != nil {
		fields = append(fields, zap.String(string(KeyForRequestID), ctx.Value(KeyForRequestID).(string)))
	}
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(KeyForRequestID) != nil {
		fields = append(fields, zap.String(string(KeyForRequestID), ctx.Value(KeyForRequestID).(string)))
	}
	l.l.Error(msg, fields...)
}
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(KeyForRequestID) != nil {
		fields = append(fields, zap.String(string(KeyForRequestID), ctx.Value(KeyForRequestID).(string)))
	}
	l.l.Fatal(msg, fields...)
}
