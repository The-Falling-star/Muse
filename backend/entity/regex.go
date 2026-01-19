package entity

import (
	"time"
)

// RegexRule 正则规则表实体
// PresetID 为 0 表示全局正则规则，不依附于任何预设
// CharacterID 为 0 表示非角色范围正则，否则关联到特定角色（Scoped Scripts）
type RegexRule struct {
	ID                      int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PresetID                int       `gorm:"column:preset_id;default:0;index:idx_preset_id" json:"presetId"`
	CharacterID             int       `gorm:"column:character_id;default:0;index:idx_character_id" json:"characterId"`
	Name                    string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	FindPattern             string    `gorm:"column:find_pattern;type:text;not null" json:"findPattern"`
	ReplacePattern          string    `gorm:"column:replace_pattern;type:text" json:"replacePattern,omitempty"`
	IsEnabled               bool      `gorm:"column:is_enabled;not null;default:true" json:"isEnabled"`
	RunOnEdit               bool      `gorm:"column:run_on_edit;not null;default:true" json:"runOnEdit"`
	SubstituteRegex         bool      `gorm:"column:substitute_regex;not null;default:true" json:"substituteRegex"`
	MinDepth                int       `gorm:"column:min_depth" json:"minDepth,omitempty"`
	MaxDepth                int       `gorm:"column:max_depth" json:"maxDepth,omitempty"`
	AffectFlagsUserInput    bool      `gorm:"column:affect_flags_user_input;not null;default:0" json:"affectFlagsUserInput"`
	AffectFlagsAIOutput     bool      `gorm:"column:affect_flags_ai_output;not null;default:false" json:"affectFlagsAIOutput"`
	AffectFlagsSlashCommand bool      `gorm:"column:affect_flags_slash_command;not null;default:false" json:"affectFlagsSlashCommand"`
	AffectFlagsWorldInfo    bool      `gorm:"column:affect_flags_world_info;not null;default:false" json:"affectFlagsWorldInfo"`
	AffectFlagsPrompt       bool      `gorm:"column:affect_flags_prompt;not null;default:false" json:"affectFlagsPrompt"`
	SortOrder               int       `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sortOrder"`
	CreatedAt               time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt               time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回表名
func (RegexRule) TableName() string {
	return "regex_rules"
}
