// Package entity 定义数据库模型实体类
package entity

import (
	"time"

	pb "github.com/ling/muse/gen/muse"
)

// User 用户表实体
type User struct {
	ID              int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                      // 用户唯一标识
	Username        string         `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_username" json:"username"` // 用户名，唯一
	PasswordHash    string         `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`                          // 密码哈希值，不返回给前端
	ActivePersonaID int            `gorm:"column:active_persona_id" json:"activePersonaId"`                                   // 当前激活的人设ID
	ActivePresetID  int            `gorm:"column:active_preset_id" json:"activePresetId"`                                     // 当前激活的预设ID
	Theme           pb.Theme       `gorm:"column:theme;not null;default:0" json:"theme"`                                      // 界面主题：0-自动，1-浅色，2-深色
	Language        string         `gorm:"column:language;type:varchar(16);not null;default:'zh-CN'" json:"language"`         // 界面语言，如zh-CN、en-US
	SendOnEnter     bool           `gorm:"column:send_on_enter;not null;default:true" json:"sendOnEnter"`                     // 是否按回车发送消息
	ShowTimestamps  bool           `gorm:"column:show_timestamps;not null;default:true" json:"showTimestamps"`                // 是否显示消息时间戳
	Provider        pb.APIProvider `gorm:"column:provider;type:tinyint unsigned;not null;default:0" json:"provider"`          // API服务提供商
	Model           string         `gorm:"column:model;type:varchar(128);not null;default:''" json:"model"`                   // 使用的模型名称
	BaseURL         string         `gorm:"column:base_url;type:varchar(512);not null;default:''" json:"baseUrl"`              // API基础URL（用于代理或自定义端点）
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                                 // 创建时间
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                                 // 更新时间

	// 关联关系
	ActivePersona Persona `gorm:"foreignKey:ActivePersonaID" json:"activePersona,omitempty"` // 当前激活的人设
	ActivePreset  Preset  `gorm:"foreignKey:ActivePresetID" json:"activePreset,omitempty"`   // 当前激活的预设
}

// TableName 返回表名
func (User) TableName() string {
	return "users"
}

// Persona 用户人设表实体
// 人设是用户在角色扮演中扮演的角色身份，类似于SillyTavern中的Persona功能
type Persona struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`            // 人设唯一标识
	UserID      int       `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"` // 所属用户ID
	Name        string    `gorm:"column:name;type:varchar(128);not null" json:"name"`      // 人设名称
	Avatar      string    `gorm:"column:avatar;type:text" json:"avatar"`                   // 人设头像，Base64编码或URL
	Description string    `gorm:"column:description;type:text" json:"description"`         // 人设描述，用于向AI描述用户扮演的角色
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`       // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`       // 更新时间
}

// TableName 返回表名
func (Persona) TableName() string {
	return "personas"
}
