package entity

import (
	"time"

	pb "github.com/ling/muse/gen/muse"
)

// Preset 预设表实体
// 预设包含模型参数配置和关联的提示项、正则规则
type Preset struct {
	ID               int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                             // 预设ID，主键自增
	UserID           int       `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`                                  // 所属用户ID
	Name             string    `gorm:"column:name;type:varchar(128);not null" json:"name"`                                       // 预设名称
	Temperature      float32   `gorm:"column:temperature;type:decimal(3,2);not null;default:1.00" json:"temperature"`            // 温度参数，控制输出随机性，范围0-2，值越高输出越随机
	TopP             float32   `gorm:"column:top_p;type:decimal(3,2);not null;default:1.00" json:"topP"`                         // 核采样参数，控制候选token的累积概率阈值，范围0-1
	TopK             int       `gorm:"column:top_k;not null;default:0" json:"topK"`                                              // Top-K采样参数，限制候选token数量，0表示不限制
	MaxTokens        int       `gorm:"column:max_tokens;not null;default:2048" json:"maxTokens"`                                 // 最大生成token数量
	FrequencyPenalty float32   `gorm:"column:frequency_penalty;type:decimal(3,2);not null;default:0.00" json:"frequencyPenalty"` // 频率惩罚，减少重复token的出现，范围0-2
	PresencePenalty  float32   `gorm:"column:presence_penalty;type:decimal(3,2);not null;default:0.00" json:"presencePenalty"`   // 存在惩罚，鼓励模型讨论新话题，范围0-2
	CandidateCount   int       `gorm:"column:candidate_count;not null;default:1" json:"candidateCount"`                          // 候选回复数量，生成多个回复供用户选择
	Version          int       `gorm:"column:version;not null;default:1" json:"version"`                                         // 乐观锁版本号，用于并发控制
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                                        // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                                        // 更新时间

	// 关联关系
	PromptItems []PromptItem `gorm:"foreignKey:PresetID" json:"promptItems,omitempty"` // 关联的提示项列表
	RegexRules  []RegexRule  `gorm:"foreignKey:PresetID" json:"regexRules,omitempty"`  // 关联的正则规则列表
}

// TableName 返回表名
func (Preset) TableName() string {
	return "presets"
}

// PromptItem 提示项表实体
// 提示项是预设中的提示词模板，可以按角色和位置插入到对话中
type PromptItem struct {
	ID                int                     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                                // 提示项ID，主键自增
	PresetID          int                     `gorm:"column:preset_id;not null;index:idx_preset_id" json:"presetId"`                               // 所属预设ID
	Identifier        pb.PromptItemIdentifier `gorm:"column:identifier;type:tinyint unsigned;not null;default:0" json:"identifier"`                // 标识符枚举，用于唯一标识提示项
	Name              string                  `gorm:"column:name;type:varchar(128);not null" json:"name"`                                          // 显示名称
	Content           string                  `gorm:"column:content;type:text" json:"content"`                                                     // 提示词内容
	Role              pb.Role                 `gorm:"column:role;type:tinyint unsigned;not null;default:0" json:"role"`                            // 消息角色：System/User/Assistant
	IsEnabled         bool                    `gorm:"column:is_enabled;not null;default:true" json:"isEnabled"`                                    // 是否启用
	InjectionPosition pb.InjectionPosition    `gorm:"column:injection_position;type:tinyint unsigned;not null;default:0" json:"injectionPosition"` // 注入位置：Relative(相对位置)/Absolute(绝对深度)
	InjectionDepth    int                     `gorm:"column:injection_depth;not null;default:0" json:"injectionDepth"`                             // 注入深度，当InjectionPosition为Absolute时生效，表示从对话末尾往前数的位置
	ForbidOverrides   bool                    `gorm:"column:forbid_overrides;not null;default:false" json:"forbidOverrides"`                       // 禁止角色卡覆盖此提示项
	SortOrder         int                     `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sortOrder"`                  // 排序顺序，数值越小越靠前
	CreatedAt         time.Time               `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                                           // 创建时间
	UpdatedAt         time.Time               `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                                           // 更新时间
}

// TableName 返回表名
func (PromptItem) TableName() string {
	return "prompt_items"
}
