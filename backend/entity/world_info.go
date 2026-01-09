package entity

import (
	"time"

	"github.com/ling/muse/common/constants"
)

// WorldInfo 世界书表实体
type WorldInfo struct {
	ID          uint32    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      uint32    `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`
	Name        string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Description *string   `gorm:"column:description;type:text" json:"description,omitempty"`
	IsGlobal    bool      `gorm:"column:is_global;not null;default:false;index:idx_is_global" json:"isGlobal"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Entries []WorldInfoEntry `gorm:"foreignKey:WorldInfoID" json:"entries,omitempty"`
}

// TableName 返回表名
func (WorldInfo) TableName() string {
	return "world_infos"
}

// WorldInfoEntry 世界书条目表实体
type WorldInfoEntry struct {
	ID             uint32                  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	WorldInfoID    uint32                  `gorm:"column:world_info_id;not null;index:idx_world_info_id" json:"worldInfoId"`
	UID            *string                 `gorm:"column:uid;type:varchar(64)" json:"uid,omitempty"`
	KeysList       string                  `gorm:"column:keys_list;type:text;not null" json:"keysList"`
	SecondaryKeys  *string                 `gorm:"column:secondary_keys;type:text" json:"secondaryKeys,omitempty"`
	Content        string                  `gorm:"column:content;type:text;not null" json:"content"`
	Comment        *string                 `gorm:"column:comment;type:text" json:"comment,omitempty"`
	IsEnabled      bool                    `gorm:"column:is_enabled;not null;default:true" json:"isEnabled"`
	Constant       bool                    `gorm:"column:constant;not null;default:false" json:"constant"`
	Selective      bool                    `gorm:"column:selective;not null;default:false" json:"selective"`
	InsertionOrder int16                   `gorm:"column:insertion_order;not null;default:100" json:"insertionOrder"`
	Position       constants.EntryPosition `gorm:"column:position;type:tinyint unsigned;not null;default:0" json:"position"`
	Depth          uint16                  `gorm:"column:depth;not null;default:4" json:"depth"`
	SortOrder      int                     `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
	CreatedAt      time.Time               `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time               `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回表名
func (WorldInfoEntry) TableName() string {
	return "world_info_entries"
}
