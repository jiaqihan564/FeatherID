package db

import (
	"fmt"
	"log"
	"time"

	"id-service/config" // 引入配置

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
	"github.com/jmoiron/sqlx"          // 高级SQL库，封装了原生database/sql
)

var DB *sqlx.DB // 全局数据库对象，供其他模块使用

// InitMySQL 初始化MySQL数据库连接
func InitMySQL(cfg config.DBConfig) error {
	// 构建MySQL连接DSN字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	// 使用sqlx打开数据库连接
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	// 设置连接池参数，保障性能与资源利用
	db.SetMaxOpenConns(50)                  // 最大连接数，根据业务调整
	db.SetMaxIdleConns(10)                  // 最大空闲连接数
	db.SetConnMaxLifetime(30 * time.Minute) // 连接最大存活时间

	// 尝试连接，确保配置正确
	if err = db.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	DB = db
	log.Println("数据库连接初始化成功")
	return nil
}
