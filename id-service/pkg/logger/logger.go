package logger

import (
	"fmt"
	"id-service/config"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/natefinch/lumberjack"
)

var (
	log         *zap.Logger
	atomicLevel zap.AtomicLevel // 支持动态调整级别
)

func InitWithConfigByDate(cfg config.LogConfig) error {
	dateStr := time.Now().Format("2006-01-02") // 获取当天日期

	logDir := filepath.Join(cfg.LogPath, dateStr)

	// 确保目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleWriter := zapcore.AddSync(os.Stdout)

	infoFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logDir, "info.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	errorFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logDir, "error.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	atomicLevel = zap.NewAtomicLevel()
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return err
	}
	atomicLevel.SetLevel(level)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleWriter, atomicLevel),
		zapcore.NewCore(encoder, infoFileWriter, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
		})),
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
