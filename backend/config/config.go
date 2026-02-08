// Package config 定义应用配置结构和加载逻辑
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Auth       AuthConfig       `mapstructure:"auth"`
	APIEncrypt APIEncryptConfig `mapstructure:"api_encrypt"`
	StaticFile StaticFileConfig `mapstructure:"static_file"`
	LogLevel   string           `mapstructure:"log_level"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Address 返回服务器监听地址
func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`      // 数据库驱动: mysql 或 sqlite
	SQLitePath   string `mapstructure:"sqlite_path"` // SQLite 数据库文件路径
	Host         string `mapstructure:"host"`        // MySQL 主机地址
	Port         int    `mapstructure:"port"`        // MySQL 端口
	Username     string `mapstructure:"username"`    // MySQL 用户名
	Password     string `mapstructure:"password"`    // MySQL 密码
	Database     string `mapstructure:"database"`    // MySQL 数据库名
	Charset      string `mapstructure:"charset"`     // MySQL 字符集
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"` // JWT 加密密钥
	SkipAuth      bool   `mapstructure:"skip_auth"`
	AdminUsername string `mapstructure:"admin_username"`
	AdminUserId   int    `mapstructure:"admin_user_id"`
	AdminPassword string `mapstructure:"admin_password"`
}

// APIEncryptConfig API密钥加密配置
type APIEncryptConfig struct {
	EncryptionKey string `mapstructure:"encryption_key"` // 加密密钥（用于AES-256）
	AllowGetKey   bool   `mapstructure:"allow_get_key"`  // 是否允许前端获取API密钥
}

// StaticFileConfig 静态文件服务配置
type StaticFileConfig struct {
	Enabled     bool   `mapstructure:"enabled"`      // 是否启用静态文件服务
	FrontendDir string `mapstructure:"frontend_dir"` // 前端构建目录路径
}

// IsSQLite 判断是否使用SQLite
func (d *DatabaseConfig) IsSQLite() bool {
	return d.Driver == "sqlite"
}

// DSN 返回MySQL数据库连接字符串
func (d *DatabaseConfig) DSN() string {
	charset := d.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		charset,
	)
}

// globalConfig 全局配置实例
var globalConfig *Config

// Load 从配置文件加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.driver", "sqlite") // 默认使用SQLite
	viper.SetDefault("database.sqlite_path", "./data/muse.database")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("auth.jwt_secret", "your-secret-key-change-in-production") // JWT密钥
	viper.SetDefault("api_encrypt.enabled", false)                              // 默认不启用API密钥加密
	viper.SetDefault("api_encrypt.allow_get_key", true)                         // 默认允许获取API密钥
	viper.SetDefault("static_file.enabled", true)                               // 默认启用静态文件服务
	viper.SetDefault("static_file.frontend_dir", "./frontend/dist")             // 默认前端目录

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = cfg
	return cfg, nil
}

// Get 获取全局配置实例
func Get() *Config {
	return globalConfig
}

// SetGlobal 设置全局配置实例（主要用于测试）
func SetGlobal(cfg *Config) {
	globalConfig = cfg
}
