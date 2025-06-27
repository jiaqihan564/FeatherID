package logger

import (
	"id-service/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/natefinch/lumberjack"
)

var (
	log         *zap.Logger
	atomicLevel zap.AtomicLevel // 支持动态调整级别
)

// InitWithConfig 完善：区分 info.log、error.log 文件
func InitWithConfig(cfg config.LogConfig) error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleWriter := zapcore.AddSync(os.Stdout)

	infoFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogPath + "/info.log",
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	errorFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogPath + "/error.log",
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	// 动态级别控制
	atomicLevel = zap.NewAtomicLevel()
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}
	atomicLevel.SetLevel(level)

	// 核心组合：
	core := zapcore.NewTee(
		// 控制台输出，所有级别
		zapcore.NewCore(encoder, consoleWriter, atomicLevel),

		// info.log 文件，info及以上
		zapcore.NewCore(encoder, infoFileWriter, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
		})),

		// error.log 文件，error及以上
		zapcore.NewCore(encoder, errorFileWriter, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	)

	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

// parseLevel 解析日志级别字符串
func parseLevel(levelStr string) (zapcore.Level, error) {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, nil // 默认info
	}
}

// SetLevel 动态修改日志级别
func SetLevel(levelStr string) {
	level, err := parseLevel(levelStr)
	if err != nil {
		Error("无效的日志级别", zap.String("level", levelStr))
		return
	}
	atomicLevel.SetLevel(level)
	Info("日志级别已动态调整", zap.String("new_level", levelStr))
}

// Sync 刷盘
func Sync() {
	_ = log.Sync()
}

// Info 普通日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}
