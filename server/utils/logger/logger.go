package logger

import (
	"YAccount/global"

	"go.uber.org/zap"
)

func GetLogger() *zap.Logger {
	return global.Logger
}

func Debug(msg string, fields ...zap.Field) {
	global.Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	global.Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	global.Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	global.Logger.Error(msg, fields...)
}
