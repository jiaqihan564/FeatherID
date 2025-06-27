package main

import (
	"context"
	"id-service/config"
	"id-service/internal/api"
	"id-service/internal/db"
	"id-service/internal/service"
	"id-service/pkg/logger" // 引入日志包
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	logCfg := config.LogConfig{
		Level:      "info",   // 初始日志级别
		LogPath:    "./logs", // 日志文件路径
		MaxSize:    100,      // 单文件最大100MB
		MaxBackups: 5,        // 最多保留5个历史文件
		MaxAge:     30,       // 日志保留30天
		Compress:   true,     // 是否压缩旧日志
	}

	// 初始化日志
	if err := logger.InitWithConfigByDate(logCfg); err != nil {
		panic("日志系统初始化失败: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("日志系统初始化成功")

	go func() {
		for {
			now := time.Now()
			// 计算明天0点时间
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			sleep := time.Until(next)

			time.Sleep(sleep)

			logger.Info("触发日志热切换")
			if err := logger.RebuildLogger(); err != nil {
				logger.Error("日志热切换失败", zap.Error(err))
			} else {
				logger.Info("日志热切换成功")
			}
		}
	}()

	// 数据库配置
	cfg := config.DBConfig{
		Host:     "192.168.200.130",
		Port:     3306,
		User:     "root",
		Password: "mysql_F7KJNF",
		DBName:   "id_service",
	}

	// 初始化数据库
	if err := db.InitMySQL(cfg); err != nil {
		logger.Error("数据库初始化失败", zap.Error(err))
		panic(err)
	}
	logger.Info("数据库连接成功")

	idGen := service.NewGenerator()
	handler := api.NewHandler(idGen)

	http.HandleFunc("/api/v1/id", handler.GetIDHandler)
	http.HandleFunc("/api/v1/id/batch", handler.GetBatchIDHandler)
	http.HandleFunc("/api/v1/log/level", handler.SetLogLevelHandler)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("服务启动，监听端口8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("监听失败", zap.Error(err))
			panic(err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("收到退出信号，准备关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器优雅关闭失败", zap.Error(err))
	} else {
		logger.Info("服务器优雅关闭完成")
	}
}
