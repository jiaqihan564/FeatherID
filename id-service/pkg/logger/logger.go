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

// InitWithConfig 通过配置初始化日志
func InitWithConfig(cfg config.LogConfig) error {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 设置动态级别控制器
	atomicLevel = zap.NewAtomicLevel()

	// 解析配置中的日志级别
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}
	atomicLevel.SetLevel(level)

	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleWriter, atomicLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, atomicLevel),
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
