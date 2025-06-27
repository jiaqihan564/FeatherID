package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"id-service/config"
	"id-service/internal/api"
	"id-service/internal/db"
	"id-service/internal/service"
)

func main() {
	// 1. 加载配置（示例硬编码，可扩展为读取文件或环境变量）
	cfg := config.DBConfig{
		Host:     "192.168.200.131",
		Port:     3306,
		User:     "root",
		Password: "mysql_F7KJNF",
		DBName:   "id_service",
	}

	// 2. 初始化数据库连接
	if err := db.InitMySQL(cfg); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 3. 初始化ID生成器
	idGen := service.NewGenerator()

	// 4. 初始化HTTP Handler
	handler := api.NewHandler(idGen)

	// 5. 路由配置
	http.HandleFunc("/api/v1/id", handler.GetIDHandler)
	http.HandleFunc("/api/v1/id/batch", handler.GetBatchIDHandler)

	// 6. 启动HTTP服务
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// 7. 优雅关闭处理
	go func() {
		log.Println("服务启动，监听端口8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听错误: %v", err)
		}
	}()

	// 8. 监听系统信号，实现优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("服务器关闭中...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭异常: %v", err)
	}

	log.Println("服务器优雅关闭完成")
}
