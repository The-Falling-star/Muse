package entity

import (
	"time"
)

// Character 角色卡表实体
type Character struct {
	ID              uint32    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID          uint32    `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`
	Name            string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Avatar          *string   `gorm:"column:avatar;type:text" json:"avatar,omitempty"`
	Description     *string   `gorm:"column:description;type:text" json:"description,omitempty"`
	Personality     *string   `gorm:"column:personality;type:text" json:"personality,omitempty"`
	Scenario        *string   `gorm:"column:scenario;type:text" json:"scenario,omitempty"`
	FirstMessage    *string   `gorm:"column:first_message;type:text" json:"firstMessage,omitempty"`
	ExampleDialogue *string   `gorm:"column:example_dialogue;type:text" json:"exampleDialogue,omitempty"`
	CreatorNotes    *string   `gorm:"column:creator_notes;type:text" json:"creatorNotes,omitempty"`
	SystemPrompt    *string   `gorm:"column:system_prompt;type:text" json:"systemPrompt,omitempty"`
	WorldInfoID     *uint32   `gorm:"column:world_info_id;index:idx_world_info_id" json:"worldInfoId,omitempty"`
	WorldInfoBackup []byte    `gorm:"column:world_info_backup;type:mediumblob" json:"-"`
	Version         uint32    `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	WorldInfo *WorldInfo `gorm:"foreignKey:WorldInfoID" json:"worldInfo,omitempty"`
}

// TableName 返回表名
func (Character) TableName() string {
	return "characters"
}
