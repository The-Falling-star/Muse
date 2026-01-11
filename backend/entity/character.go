package entity

import (
	"time"
)

// Character 角色卡表实体
type Character struct {
	ID     int `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID int `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`

	// 用户能够编辑的内容
	Name            string `gorm:"column:name;type:varchar(128);not null" json:"name"`       // 角色名称
	Avatar          string `gorm:"column:avatar;type:text" json:"avatar"`                    // 角色头像
	Description     string `gorm:"column:description;type:text" json:"description"`          // 角色描述
	FirstMessage    string `gorm:"column:first_message;type:text" json:"firstMessage"`       // 角色开场白
	ExampleDialogue string `gorm:"column:example_dialogue;type:text" json:"exampleDialogue"` // 示例对话
	CreatorNotes    string `gorm:"column:creator_notes;type:text" json:"creatorNotes"`       // 创作者备注

	// 数据库字段内容
	WorldInfoID     int       `gorm:"column:world_info_id;index:idx_world_info_id" json:"worldInfoId"` // 关联的世界信息ID
	WorldInfoBackup []byte    `gorm:"column:world_info_backup;type:mediumblob" json:"-"`               // 世界信息备份，存储PNG数据
	Version         int       `gorm:"column:lock_version;not null;default:1" json:"lock_version"`      // 乐观锁版本号
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`               // 创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`               // 更新时间
	// 关联关系
	WorldInfo WorldInfo `gorm:"foreignKey:WorldInfoID" json:"worldInfo,omitempty"` // 关联的世界信息
}

// TableName 返回表名
func (Character) TableName() string {
	return "characters"
}
