# FeatherID - 分布式ID生成服务

## 项目概述

FeatherID 是一个基于号段模式的分布式ID生成服务，采用Go语言开发，提供高性能、高可用的唯一标识符生成能力。

### 核心特性

- **高性能**: 基于号段缓存机制，减少数据库访问频率
- **高可用**: 支持多实例部署，通过数据库事务保证ID唯一性
- **业务隔离**: 支持多业务线独立管理，避免ID冲突
- **动态配置**: 支持运行时调整日志级别，无需重启服务
- **优雅关闭**: 支持信号处理，确保服务安全停止

## 技术栈

- **开发语言**: Go 1.24.4
- **Web框架**: net/http (标准库)
- **数据库**: MySQL 8.0+
- **数据库驱动**: go-sql-driver/mysql v1.9.3
- **SQL工具**: jmoiron/sqlx v1.4.0
- **日志系统**: uber-go/zap v1.27.0
- **日志轮转**: lumberjack v2.0.0

## 项目结构

```
id-service/
├── cmd/main.go                 # 应用程序入口
├── config/config.go            # 配置结构定义
├── internal/
│   ├── api/handler.go          # HTTP接口处理器
│   ├── db/mysql.go             # 数据库连接管理
│   ├── model/segment.go        # 数据模型定义
│   └── service/
│       ├── id_generator.go     # ID生成器核心逻辑
│       └── segment_service.go  # 号段服务
└── pkg/logger/logger.go        # 日志系统
```

## API接口

### 1. 获取单个ID
```
GET /api/v1/id?biz_tag=user
```

### 2. 批量获取ID
```
GET /api/v1/id/batch?biz_tag=order&count=10
```

### 3. 动态调整日志级别
```
GET /api/v1/log/level?level=debug
```

## 数据库设计

### id_segments 表结构
```sql
CREATE TABLE `id_segments` (
  `biz_tag` varchar(64) NOT NULL COMMENT '业务标识',
  `max_id` bigint NOT NULL DEFAULT '0' COMMENT '当前最大ID值',
  `step` int NOT NULL DEFAULT '1000' COMMENT '号段步长',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`biz_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ID号段表';
```

## 部署说明

### 环境要求
- Go 1.24.4+
- MySQL 8.0+
- 内存: 至少512MB

### 快速部署

1. **安装依赖**
```bash
cd id-service
go mod download
```

2. **编译项目**
```bash
go build -o featherid cmd/main.go
```

3. **配置数据库**
```sql
CREATE DATABASE id_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE `id_segments` (
  `biz_tag` varchar(64) NOT NULL,
  `max_id` bigint NOT NULL DEFAULT '0',
  `step` int NOT NULL DEFAULT '1000',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`biz_tag`)
);
INSERT INTO id_segments (biz_tag, max_id, step) VALUES ('user', 0, 1000), ('order', 0, 1000);
```

4. **修改配置**
编辑 `cmd/main.go` 中的数据库配置

5. **启动服务**
```bash
./featherid
```

## 性能特性

- **单机QPS**: 10,000+ ID/秒
- **响应时间**: 平均 < 1ms
- **并发支持**: 支持1000+并发请求
- **内存占用**: 约50MB

## 核心机制

### ID生成流程
1. 客户端请求ID，携带业务标识(biz_tag)
2. 检查本地缓存是否有可用号段
3. 如果缓存为空或已用完，从数据库申请新号段
4. 使用数据库事务确保多实例下ID不重复
5. 返回递增的ID给客户端

### 并发安全
- 使用读写锁保护缓存访问
- 使用互斥锁保护号段操作
- 数据库层面使用FOR UPDATE锁定行

## 监控与运维

### 日志管理
- 支持debug、info、warn、error四个级别
- 按日期自动分割日志文件
- 支持日志轮转和压缩
- 支持运行时动态调整日志级别

### 健康检查
```bash
curl "http://localhost:8080/api/v1/id?biz_tag=health"
```

---

**FeatherID** - 让ID生成变得简单高效 🚀 