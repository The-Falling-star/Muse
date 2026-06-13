package entity

import (
	"time"
)

// RegexRule 正则规则表实体
// PresetID 为 0 表示全局正则规则，不依附于任何预设
// CharacterID 为 0 表示非角色范围正则，否则关联到特定角色（Scoped Scripts）
type RegexRule struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                        // 主键ID
	UserID      int    `gorm:"column:user_id;not null;default:0;index:idx_regex_rules_user_id" json:"userId"`       // 用户ID
	PresetID    int    `gorm:"column:preset_id;default:0;index:idx_regex_rules_preset_id" json:"presetId"`          // 关联的预设ID，0表示全局正则
	CharacterID int    `gorm:"column:character_id;default:0;index:idx_regex_rules_character_id" json:"characterId"` // 关联的角色ID，0表示非角色范围正则
	Name        string `gorm:"column:name;type:varchar(128);not null" json:"name"`                                  // 规则名称

	// 正则表达式配置
	FindPattern    string `gorm:"column:find_pattern;type:text;not null" json:"findPattern"`        // 查找模式（正则表达式）
	ReplacePattern string `gorm:"column:replace_pattern;type:text" json:"replacePattern,omitempty"` // 替换模式

	// 规则开关和行为
	IsEnabled       bool `gorm:"column:is_enabled;not null;default:true" json:"isEnabled"`             // 是否启用
	RunOnEdit       bool `gorm:"column:run_on_edit;not null;default:true" json:"runOnEdit"`            // 编辑消息时是否运行
	SubstituteRegex bool `gorm:"column:substitute_regex;not null;default:true" json:"substituteRegex"` // 是否使用正则替换（false则使用字面量替换）

	// 深度范围限制
	MinDepth int `gorm:"column:min_depth;not null;default:0" json:"minDepth,omitempty"`          // 最小深度（从对话末尾开始计算）
	MaxDepth int `gorm:"column:max_depth;not null;default:2147483647" json:"maxDepth,omitempty"` // 最大深度（从对话末尾开始计算）

	// 作用范围标志
	AffectFlagsUserInput    bool `gorm:"column:affect_flags_user_input;not null;default:0" json:"affectFlagsUserInput"`           // 是否作用于用户输入
	AffectFlagsAIOutput     bool `gorm:"column:affect_flags_ai_output;not null;default:false" json:"affectFlagsAIOutput"`         // 是否作用于AI输出
	AffectFlagsSlashCommand bool `gorm:"column:affect_flags_slash_command;not null;default:false" json:"affectFlagsSlashCommand"` // 是否作用于斜杠命令
	AffectFlagsWorldInfo    bool `gorm:"column:affect_flags_world_info;not null;default:false" json:"affectFlagsWorldInfo"`       // 是否作用于世界书
	AffectFlagsPrompt       bool `gorm:"column:affect_flags_prompt;not null;default:false" json:"affectFlagsPrompt"`              // 是否作用于提示词

	// 排序和时间
	SortOrder int       `gorm:"column:sort_order;not null;default:0;index:idx_regex_rules_sort_order" json:"sortOrder"` // 排序顺序，数值越小越靠前
	Version   int       `gorm:"column:version;not null;default:1" json:"version"`                                       // 乐观锁版本号，用于缓存校验
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                                      // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                                      // 更新时间
}

// TableName 返回表名
func (RegexRule) TableName() string {
	return "regex_rules"
}
