package entity

import (
	"time"

	pb "github.com/ling/muse/gen/muse"
)

// APIConfig API配置表实体
// 用于存储用户配置的各类大模型API信息，如OpenAI、Claude、Gemini等
type APIConfig struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                // 主键ID
	UserID    int            `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`                     // 所属用户ID
	Provider  pb.APIProvider `gorm:"column:provider;type:tinyint unsigned;not null" json:"provider"`              // API提供商类型（OpenAI/Claude/Gemini）
	APIKey    string         `gorm:"column:api_key;type:varchar(512);not null" json:"-"`                          // API密钥，不返回给前端
	IsActive  bool           `gorm:"column:is_active;not null;default:false;index:idx_is_active" json:"isActive"` // 是否为当前激活的配置
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                           // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                           // 更新时间
}

// TableName 返回表名
func (APIConfig) TableName() string {
	return "api_configs"
}
