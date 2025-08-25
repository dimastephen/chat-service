package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var authLogger *zap.Logger

func Init(core zapcore.Core, opts ...zap.Option) {
	authLogger = zap.New(core, opts...)
}

func Debug(msg string, fields ...zap.Field) {
	authLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	authLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	authLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	authLogger.Error(msg, fields...)
}
