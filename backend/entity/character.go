package entity

import (
	"time"
)

// Character 角色卡表实体
// 存储AI角色的基本信息和设定，用于角色扮演对话
type Character struct {
	ID     int `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                       // 主键ID
	UserID int `gorm:"column:user_id;not null;index:idx_characters_user_id" json:"userId"` // 所属用户ID

	// 用户能够编辑的内容
	Name            string   `gorm:"column:name;type:varchar(128);not null" json:"name"`                       // 角色名称
	Avatar          string   `gorm:"column:avatar;type:varchar(1024)" json:"avatar"`                           // 角色头像路径
	Description     string   `gorm:"column:description;type:mediumtext" json:"description"`                    // 角色描述/人设，包含角色的性格、背景等信息
	FirstMessage    []string `gorm:"column:first_message;type:json;serializer:json" json:"firstMessage"`       // 角色开场白，对话开始时角色发送的第一条消息
	ExampleDialogue []string `gorm:"column:example_dialogue;type:json;serializer:json" json:"exampleDialogue"` // 示例对话，用于指导AI如何扮演该角色
	CreatorNotes    string   `gorm:"column:creator_notes;type:text" json:"creatorNotes"`                       // 创作者备注，角色卡创建者留下的说明

	// 数据库字段内容
	WorldInfoID     int       `gorm:"column:world_info_id;index:idx_characters_world_info_id" json:"worldInfoId"` // 关联的世界书ID，0表示无关联
	WorldInfoBackup []byte    `gorm:"column:world_info_backup;type:mediumblob" json:"-"`                          // 世界书备份数据，存储导入时角色卡自带的世界书原始JSON
	Version         int       `gorm:"column:lock_version;not null;default:1" json:"lock_version"`                 // 乐观锁版本号，用于并发更新控制
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                          // 创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                          // 更新时间

	// 关联关系
	WorldInfo  WorldInfo   `gorm:"foreignKey:WorldInfoID" json:"worldInfo,omitempty"` // 关联的世界书对象
	RegexRules []RegexRule `gorm:"foreignKey:CharacterID" json:"regexRules,omitempty"`
}

// TableName 返回表名
func (Character) TableName() string {
	return "characters"
}
