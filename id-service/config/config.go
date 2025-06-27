package config

// DBConfig 数据库配置结构体
type DBConfig struct {
	Host     string // 数据库主机地址
	Port     int    // 数据库端口
	User     string // 数据库用户名
	Password string // 数据库密码
	DBName   string // 数据库名称
}
