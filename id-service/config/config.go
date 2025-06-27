package config

// DBConfig 数据库配置结构体
type DBConfig struct {
	Host     string // 数据库主机地址
	Port     int    // 数据库端口
	User     string // 数据库用户名
	Password string // 数据库密码
	DBName   string // 数据库名称
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string // 日志级别: debug、info、warn、error
	LogPath    string // 日志文件路径
	MaxSize    int    // 单文件最大MB
	MaxBackups int    // 保留旧文件数量
	MaxAge     int    // 文件保留天数
	Compress   bool   // 是否压缩
}
