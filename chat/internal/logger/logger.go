package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var chatLogger *zap.Logger

func Init(core zapcore.Core, opts ...zap.Option) {
	chatLogger = zap.New(core, opts...)
}

func Debug(msg string, fields ...zap.Field) {
	chatLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	chatLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	chatLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	chatLogger.Error(msg, fields...)
}
