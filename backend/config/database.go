package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ling/muse/entity"
)

// db 全局数据库连接实例
var db *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *DatabaseConfig) error {
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.IsSQLite() {
		db, err = initSQLite(cfg, gormConfig)
	} else {
		db, err = initMySQL(cfg, gormConfig)
	}

	return err
}

// initSQLite 初始化SQLite数据库
func initSQLite(cfg *DatabaseConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	// 确保数据目录存在
	dbDir := filepath.Dir(cfg.SQLitePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	database, err := gorm.Open(sqlite.Open(cfg.SQLitePath), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接SQLite数据库失败: %w", err)
	}

	// SQLite 开启外键约束
	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("开启外键约束失败: %w", err)
	}

	// SQLite 模式自动建表
	if err := autoMigrate(database); err != nil {
		return nil, fmt.Errorf("自动建表失败: %w", err)
	}

	return database, nil
}

// initMySQL 初始化MySQL数据库
func initMySQL(cfg *DatabaseConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	database, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接MySQL数据库失败: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return database, nil
}

// autoMigrate 自动迁移数据库表结构（仅SQLite模式使用）
func autoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&entity.User{},
		&entity.Persona{},
		&entity.UserSetting{},
		&entity.Character{},
		&entity.ChatSession{},
		&entity.Message{},
		&entity.MessageSwipe{},
		&entity.Preset{},
		&entity.PromptItem{},
		&entity.RegexRule{},
		&entity.WorldInfo{},
		&entity.WorldInfoEntry{},
		&entity.APIConfig{},
	)
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return db
}
