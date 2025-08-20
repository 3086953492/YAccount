package initialize

import (
	"YAccount/global"
	"YAccount/utils/logger"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() error {
	cfg := global.Cfg.Log

	// 确保日志目录存在
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return err
	}

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var fileWriter zapcore.WriteSyncer

	// 根据配置选择日志轮转方式
	if cfg.RotateDaily {
		// 使用按日期分割的写入器
		dailyWriter := logger.NewDailyRotateWriter(
			cfg.Filename,
			cfg.MaxSize,
			cfg.MaxBackups,
			cfg.MaxAge,
			cfg.Compress,
		)
		fileWriter = zapcore.AddSync(dailyWriter)
	} else {
		// 使用原来的lumberjack按大小分割
		lumberJackLogger := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		fileWriter = zapcore.AddSync(lumberJackLogger)
	}

	// 创建核心
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			level,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
	)

	global.Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}
