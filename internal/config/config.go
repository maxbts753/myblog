package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// AppConfig 全局配置实例
// 存储应用程序的所有配置信息
var AppConfig *Config

// Config 应用程序主配置结构体
// 包含应用程序的所有配置部分
// 使用mapstructure标签与viper配合实现配置文件映射
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // 服务器相关配置
	Database DatabaseConfig `mapstructure:"database"` // 数据库相关配置
}
// ServerConfig 服务器配置结构体
// 包含HTTP服务器的配置信息
type ServerConfig struct {
	Port string `mapstructure:"port"` // 服务器监听端口
}

// DatabaseConfig 数据库配置结构体
// 包含数据库连接的所有参数
// 注意：在当前简化实现中，部分字段可能未被实际使用
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

// Init 初始化配置
// 此函数负责：
// 1. 设置viper配置文件名和类型
// 2. 添加配置文件搜索路径
// 3. 设置默认配置值
// 4. 读取配置文件
// 5. 将配置解析到AppConfig全局变量
// 注意：
//   如果配置文件读取失败，函数会直接panic，终止程序启动
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	
	// 设置默认配置值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "root")
	viper.SetDefault("database.dbname", "blog_db")

	// 尝试读取配置文件，如果失败则打印警告但不panic
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: Failed to read config file, using default values")
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		panic(err)
	}
}
