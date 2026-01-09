package entity

import (
	"time"

	"github.com/ling/muse/common/constants"
)

// APIConfig API配置表实体
type APIConfig struct {
	ID        uint32                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint32                `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`
	Name      string                `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Provider  constants.APIProvider `gorm:"column:provider;type:tinyint unsigned;not null" json:"provider"`
	APIKey    string                `gorm:"column:api_key;type:varchar(512);not null" json:"-"`
	BaseURL   *string               `gorm:"column:base_url;type:varchar(512)" json:"baseUrl,omitempty"`
	Model     *string               `gorm:"column:model;type:varchar(128)" json:"model,omitempty"`
	IsActive  bool                  `gorm:"column:is_active;not null;default:false;index:idx_is_active" json:"isActive"`
	CreatedAt time.Time             `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time             `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回表名
func (APIConfig) TableName() string {
	return "api_configs"
}
