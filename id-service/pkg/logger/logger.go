package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger // 全局日志对象

// Init 初始化日志系统
func Init() error {
	var err error
	// 这里使用 zap.NewProduction() 生产环境推荐配置，日志结构清晰
	log, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

// Sync 手动刷盘，确保日志完整写入
func Sync() {
	_ = log.Sync()
}

// Info 普通信息日志
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
