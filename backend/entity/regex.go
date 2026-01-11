package entity

import (
	"time"

	"github.com/ling/muse/common/constants"
)

// RegexRule 正则规则表实体
type RegexRule struct {
	ID              int                        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PresetID        int                        `gorm:"column:preset_id;not null;index:idx_preset_id" json:"presetId"`
	Name            string                     `gorm:"column:name;type:varchar(128);not null" json:"name"`
	FindPattern     string                     `gorm:"column:find_pattern;type:text;not null" json:"findPattern"`
	ReplacePattern  string                     `gorm:"column:replace_pattern;type:text" json:"replacePattern,omitempty"`
	IsEnabled       bool                       `gorm:"column:is_enabled;not null;default:true" json:"isEnabled"`
	RunOnEdit       bool                       `gorm:"column:run_on_edit;not null;default:true" json:"runOnEdit"`
	SubstituteRegex bool                       `gorm:"column:substitute_regex;not null;default:true" json:"substituteRegex"`
	MinDepth        int                        `gorm:"column:min_depth" json:"minDepth,omitempty"`
	MaxDepth        int                        `gorm:"column:max_depth" json:"maxDepth,omitempty"`
	AffectFlags     constants.RegexAffectFlags `gorm:"column:affect_flags;type:int unsigned;not null;default:0" json:"affectFlags"`
	SortOrder       int                        `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sortOrder"`
	CreatedAt       time.Time                  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time                  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回表名
func (RegexRule) TableName() string {
	return "regex_rules"
}
