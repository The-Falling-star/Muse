package entity

import (
	"time"

	pb "github.com/ling/muse/gen/muse"
)

// Preset 预设表实体
type Preset struct {
	ID               int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID           int       `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`
	Name             string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Temperature      float64   `gorm:"column:temperature;type:decimal(3,2);not null;default:1.00" json:"temperature"`
	TopP             float64   `gorm:"column:top_p;type:decimal(3,2);not null;default:1.00" json:"topP"`
	TopK             int       `gorm:"column:top_k;not null;default:0" json:"topK"`
	MaxTokens        int       `gorm:"column:max_tokens;not null;default:2048" json:"maxTokens"`
	FrequencyPenalty float64   `gorm:"column:frequency_penalty;type:decimal(3,2);not null;default:0.00" json:"frequencyPenalty"`
	PresencePenalty  float64   `gorm:"column:presence_penalty;type:decimal(3,2);not null;default:0.00" json:"presencePenalty"`
	Version          int       `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	PromptItems []PromptItem `gorm:"foreignKey:PresetID" json:"promptItems,omitempty"`
	RegexRules  []RegexRule  `gorm:"foreignKey:PresetID" json:"regexRules,omitempty"`
}

// TableName 返回表名
func (Preset) TableName() string {
	return "presets"
}

// PromptItem 提示项表实体
type PromptItem struct {
	ID                int                  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PresetID          int                  `gorm:"column:preset_id;not null;index:idx_preset_id" json:"presetId"`
	Identifier        string               `gorm:"column:identifier;type:varchar(64);not null" json:"identifier"`
	Name              string               `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Content           string               `gorm:"column:content;type:text" json:"content"`
	Role              pb.Role              `gorm:"column:role;type:tinyint unsigned;not null;default:0" json:"role"`
	IsEnabled         bool                 `gorm:"column:is_enabled;not null;default:true" json:"isEnabled"`
	InjectionPosition pb.InjectionPosition `gorm:"column:injection_position;type:tinyint unsigned;not null;default:0" json:"injectionPosition"`
	InjectionDepth    int                  `gorm:"column:injection_depth;not null;default:0" json:"injectionDepth"`
	ForbidOverrides   bool                 `gorm:"column:forbid_overrides;not null;default:false" json:"forbidOverrides"`
	SortOrder         int                  `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sortOrder"`
	CreatedAt         time.Time            `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time            `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回表名
func (PromptItem) TableName() string {
	return "prompt_items"
}
