package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/natefinch/lumberjack"
)

var log *zap.Logger // 全局日志对象

// Init 初始化日志，输出到控制台和文件
func Init() error {
	// 配置 lumberjack 文件切割
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    100,              // 单个文件最大100MB
		MaxBackups: 5,                // 最多保留5个旧文件
		MaxAge:     30,               // 文件保留30天
		Compress:   true,             // 压缩旧日志文件
	})

	// 编码配置，结构化JSON输出
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts" // 时间字段名
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 设置日志级别
	level := zap.InfoLevel

	// 构建日志核心
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleWriter, level),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileWriter, level),
	)

	// 最终初始化日志对象
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
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
