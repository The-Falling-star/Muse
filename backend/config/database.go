package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite" // pure Go SQLite driver
	"github.com/ling/muse/common/crypto"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
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
		Logger:                                   logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	}

	if cfg.IsSQLite() {
		db, err = initSQLite(cfg, gormConfig)
	} else {
		db, err = initMySQL(cfg, gormConfig)
	}
	if err != nil {
		return err
	}
	var admin entity.User
	result := db.Where("id = ?", Get().Auth.AdminUserId).First(&admin)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Infof("管理员记录不存在，将创建管理员")
		// 记录不存在，创建管理员
		admin.ID = Get().Auth.AdminUserId
		admin.Username = Get().Auth.AdminUsername
		hashPassword, _ := crypto.Md5HashStr(Get().Auth.AdminPassword)
		admin.PasswordHash = hashPassword
		if err = db.Create(&admin).Error; err != nil {
			return fmt.Errorf("创建管理员失败: %v", err)
		}
	}

	return nil
}

// initSQLite 初始化SQLite数据库
func initSQLite(cfg *DatabaseConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	// 确保数据目录存在
	dbDir := filepath.Dir(cfg.SQLitePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %v", err)
	}

	database, err := gorm.Open(sqlite.Open(cfg.SQLitePath), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
	}

	// SQLite 模式自动建表
	if err = autoMigrate(database); err != nil {
		return nil, fmt.Errorf("自动建表失败: %v", err)
	}

	return database, nil
}

// initMySQL 初始化MySQL数据库
func initMySQL(cfg *DatabaseConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	// 先尝试连接指定的数据库
	database, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	if err != nil {
		// 如果是数据库不存在的错误，尝试创建数据库
		if !strings.Contains(err.Error(), "Unknown database") {
			return nil, fmt.Errorf("连接MySQL数据库失败: %v", err)
		}

		log.Infof("数据库 %s 不存在，尝试自动创建", cfg.Database)
		if createErr := createDatabase(cfg, gormConfig); createErr != nil {
			return nil, fmt.Errorf("创建数据库失败: %v", createErr)
		}

		// 重新连接
		database, err = gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("连接MySQL数据库失败: %v", err)
		}
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// MySQL 模式也自动建表
	if err = autoMigrate(database); err != nil {
		return nil, fmt.Errorf("自动建表失败: %v", err)
	}

	return database, nil
}

// createDatabase 创建数据库
func createDatabase(cfg *DatabaseConfig, gormConfig *gorm.Config) error {
	// 连接到 MySQL 服务器（不指定数据库）
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Charset,
	)
	if cfg.Charset == "" {
		dsnWithoutDB = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
		)
	}

	database, err := gorm.Open(mysql.Open(dsnWithoutDB), gormConfig)
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %v", err)
	}
	defer func() {
		if sqlDB, dbErr := database.DB(); dbErr == nil {
			_ = sqlDB.Close()
		}
	}()

	// 创建数据库
	charset := cfg.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	createSQL := fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET %s COLLATE %s_general_ci", cfg.Database, charset, charset)
	if err = database.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("执行创建数据库语句失败: %v", err)
	}

	log.Infof("数据库 %s 创建成功", cfg.Database)
	return nil
}

// autoMigrate 自动迁移数据库表结构（仅SQLite模式使用）
func autoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&entity.User{},
		&entity.Persona{},
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

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}
