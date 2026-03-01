package entity

import (
	"time"

	pb "github.com/ling/muse/gen/muse"
)

// WorldInfo 世界书表实体
// 世界书是一组条目的集合，用于在对话中根据关键词触发插入额外的上下文信息
type WorldInfo struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                // 世界书唯一标识
	UserID      int       `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`                     // 所属用户ID
	Name        string    `gorm:"column:name;type:varchar(128);not null;default:''" json:"name"`               // 世界书名称
	Description string    `gorm:"column:description;type:text;not null" json:"description,omitempty"`          // 世界书描述
	IsGlobal    bool      `gorm:"column:is_global;not null;default:false;index:idx_is_global" json:"isGlobal"` // 是否为全局世界书，全局世界书对所有角色生效
	Version     int       `gorm:"column:version;not null;default:1" json:"version"`                            // 乐观锁版本号，用于缓存校验
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                           // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                           // 更新时间

	// 关联关系
	Entries []WorldInfoEntry `gorm:"foreignKey:WorldInfoID" json:"entries,omitempty"` // 世界书条目列表
}

// TableName 返回表名
func (WorldInfo) TableName() string {
	return "world_infos"
}

// WorldInfoEntry 世界书条目表实体
// 条目包含触发关键词和对应的内容，当对话中出现关键词时会将内容插入到提示词中
type WorldInfoEntry struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                             // 条目唯一标识
	WorldInfoID int    `gorm:"column:world_info_id;not null;index:idx_world_info_id" json:"worldInfoId"` // 所属世界书ID
	UID         string `gorm:"column:uid;type:varchar(64)" json:"uid"`                                   // 条目唯一标识符（用于导入导出）

	// 用户可编辑字段
	Keys           []string         `gorm:"column:keys;type:json;not null;serializer:json" json:"keys"`                    // 主关键词列表，逗号分隔，匹配任一关键词即可触发
	SecondaryKeys  []string         `gorm:"column:secondary_keys;type:json;not null;serializer:json" json:"secondaryKeys"` // 次要关键词列表，逗号分隔，选择性匹配时需同时匹配主次关键词
	Content        string           `gorm:"column:content;type:mediumtext;not null" json:"content"`                        // 条目内容，触发时插入到提示词中
	Comment        string           `gorm:"column:comment;type:text;not null" json:"comment"`                              // 条目备注，仅供用户参考，不会插入到提示词
	IsEnabled      bool             `gorm:"column:is_enabled;not null;default:true" json:"isEnabled"`                      // 是否启用该条目
	Constant       bool             `gorm:"column:constant;not null;default:false" json:"constant"`                        // 是否常驻，常驻条目始终插入，无需关键词触发
	Selective      bool             `gorm:"column:selective;not null;default:false" json:"selective"`                      // 是否选择性匹配，启用时需同时匹配主关键词和次要关键词
	InsertionOrder int              `gorm:"column:insertion_order;not null;default:100" json:"insertionOrder"`             // 插入顺序，数字越小越先插入（同一深度时）
	Position       pb.EntryPosition `gorm:"column:position;type:tinyint unsigned;not null;default:0" json:"position"`      // 插入位置：角色描述前/后、示例对话前/后、指定深度
	Depth          int              `gorm:"column:depth;not null;default:4" json:"depth"`                                  // 插入深度，0表示最新消息，数字越大离最新消息越远
	Role           pb.Role          `gorm:"column:role;type:tinyint unsigned;not null;default:0" json:"role"`              // 消息角色类型，用于区分该条目的作用角色
	SortOrder      int              `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`                         // 显示排序顺序

	// 数据库字段
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"` // 更新时间
}

// TableName 返回表名
func (WorldInfoEntry) TableName() string {
	return "world_info_entries"
}
