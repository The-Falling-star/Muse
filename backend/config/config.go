// Package config 定义应用配置结构和加载逻辑
package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "config.yaml"

// Config 应用配置结构
type Config struct {
	Server          ServerConfig        `yaml:"server"`
	Database        DatabaseConfig      `yaml:"database"`
	Auth            AuthConfig          `yaml:"auth"`
	APIEncrypt      APIEncryptConfig    `yaml:"api_encrypt"`
	StaticFile      StaticFileConfig    `yaml:"static_file"`
	File            FileConfig          `yaml:"file"`
	LogLevel        string              `yaml:"log_level" env:"LOG_LEVEL"`
	Chat            ChatConfig          `yaml:"chat"`
	CandidateModels map[string][]string `yaml:"candidate_models"`
	Cache           CacheConfig         `yaml:"cache"`
}

// ChatConfig 聊天配置
type ChatConfig struct {
	EnableCache bool `yaml:"enable_cache" env:"CHAT_ENABLE_CACHE"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST"`
	Port int    `yaml:"port" env:"SERVER_PORT"`
}

// Address 返回服务器监听地址
func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `yaml:"driver" env:"DB_DRIVER"`           // 数据库驱动: mysql 或 sqlite
	SQLitePath   string `yaml:"sqlite_path" env:"DB_SQLITE_PATH"` // SQLite 数据库文件路径
	Host         string `yaml:"host" env:"DB_HOST"`               // MySQL 主机地址
	Port         int    `yaml:"port" env:"DB_PORT"`               // MySQL 端口
	Username     string `yaml:"username" env:"DB_USERNAME"`       // MySQL 用户名
	Password     string `yaml:"password" env:"DB_PASSWORD"`       // MySQL 密码
	Database     string `yaml:"database" env:"DB_DATABASE"`       // MySQL 数据库名
	Charset      string `yaml:"charset" env:"DB_CHARSET"`         // MySQL 字符集
	MaxIdleConns int    `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns int    `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret" env:"AUTH_JWT_SECRET"` // JWT 加密密钥
	SkipAuth      bool   `yaml:"skip_auth" env:"AUTH_SKIP_AUTH"`
	AdminUsername string `yaml:"admin_username" env:"AUTH_ADMIN_USERNAME"`
	AdminUserId   int    `yaml:"admin_user_id" env:"AUTH_ADMIN_USER_ID"`
	AdminPassword string `yaml:"admin_password" env:"AUTH_ADMIN_PASSWORD"`
}

// APIEncryptConfig API密钥加密配置
type APIEncryptConfig struct {
	EncryptionKey string `yaml:"encryption_key" env:"API_ENCRYPT_KEY"`  // 加密密钥（用于AES-256）
	AllowGetKey   bool   `yaml:"allow_get_key" env:"API_ALLOW_GET_KEY"` // 是否允许前端获取API密钥
}

// StaticFileConfig 静态文件服务配置
type StaticFileConfig struct {
	Enabled     bool   `yaml:"enabled" env:"STATIC_FILE_ENABLED"`           // 是否启用静态文件服务
	FrontendDir string `yaml:"frontend_dir" env:"STATIC_FILE_FRONTEND_DIR"` // 前端构建目录路径
}

// FileConfig 文件服务配置
type FileConfig struct {
	UploadPath     string `yaml:"upload_path" env:"FILE_UPLOAD_PATH"`           // 上传文件根目录
	MaxUploadSize  int    `yaml:"max_upload_size" env:"FILE_MAX_UPLOAD_SIZE"`   // 最大上传文件大小（MB）
	CacheExpireMin int    `yaml:"cache_expire_min" env:"FILE_CACHE_EXPIRE_MIN"` // 前端缓存过期时间（分钟）
}

// CacheConfig 缓存配置
type CacheConfig struct {
	LRUMaxSize int `yaml:"lru_max_size"` // LRU缓存的最大元素个数
	LRUTTL     int `yaml:"lru_ttl"`      // LRU缓存的最大时间, 单位秒
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
	// 1. 初始化默认值 (替代之前的 viper.SetDefault)
	cfg := defaultConfig
	// 2. 读取并解析配置文件 (YAML 库会自动处理大小写敏感的 Map)
	if configPath == "" {
		configPath = DefaultConfigPath
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 如果明确指定了路径但文件不存在，则报错
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	// 3. 环境变量覆盖 (优先级最高)
	if err = env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("解析环境变量失败: %w", err)
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

// defaultConfig 默认配置
var defaultConfig = &Config{
	Server: ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
	},
	Database: DatabaseConfig{
		Driver:       "sqlite",
		SQLitePath:   "./data/muse.database",
		Host:         "localhost",
		Port:         3306,
		Charset:      "utf8mb4",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	},
	Auth: AuthConfig{
		JWTSecret: "your-secret-key-change-in-production",
	},
	APIEncrypt: APIEncryptConfig{
		AllowGetKey: true,
	},
	StaticFile: StaticFileConfig{
		Enabled:     true,
		FrontendDir: "./frontend/dist",
	},
	File: FileConfig{
		UploadPath:     "./data/uploads",
		MaxUploadSize:  10,
		CacheExpireMin: 1440,
	},
	LogLevel: "info",
	Cache: CacheConfig{
		LRUMaxSize: 100,
		LRUTTL:     300,
	},
}
