// Package model 定义数据库模型实体类
package entity

import (
	"time"
)

// User 用户表实体
type User struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username        string    `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_username" json:"username"`
	PasswordHash    string    `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Email           string    `gorm:"column:email;type:varchar(255)" json:"email"`
	ActivePersonaID int       `gorm:"column:active_persona_id" json:"activePersonaId"`
	ActivePresetID  int       `gorm:"column:active_preset_id" json:"activePresetId"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	ActivePersona Persona `gorm:"foreignKey:ActivePersonaID" json:"activePersona,omitempty"`
	ActivePreset  Preset  `gorm:"foreignKey:ActivePresetID" json:"activePreset,omitempty"`
}

// TableName 返回表名
func (User) TableName() string {
	return "users"
}

// Persona 用户人设表实体
type Persona struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`
	Name        string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Avatar      string    `gorm:"column:avatar;type:text" json:"avatar"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回表名
func (Persona) TableName() string {
	return "personas"
}

// UserSetting 用户设置表实体
type UserSetting struct {
	ID             int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID         int       `gorm:"column:user_id;not null;uniqueIndex:uk_user_id" json:"userId"`
	Theme          int       `gorm:"column:theme;not null;default:0" json:"theme"`
	Language       string    `gorm:"column:language;type:varchar(16);not null;default:'zh-CN'" json:"language"`
	SendOnEnter    bool      `gorm:"column:send_on_enter;not null;default:true" json:"sendOnEnter"`
	ShowTimestamps bool      `gorm:"column:show_timestamps;not null;default:true" json:"showTimestamps"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回表名
func (UserSetting) TableName() string {
	return "user_settings"
}
