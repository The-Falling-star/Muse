package entity

import (
	"time"

	"github.com/ling/muse/common/constants"
)

// ChatSession 聊天会话表实体
type ChatSession struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`
	CharacterID int       `gorm:"column:character_id;not null;index:idx_character_id" json:"characterId"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name,omitempty"`
	Version     int       `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Character *Character `gorm:"foreignKey:CharacterID" json:"character,omitempty"`
	Messages  []Message  `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
}

// TableName 返回表名
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// Message 消息表实体
type Message struct {
	ID               int                   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID        int                   `gorm:"column:session_id;not null;index:idx_session_id" json:"sessionId"`
	Role             constants.MessageRole `gorm:"column:role;type:tinyint unsigned;not null" json:"role"`
	ActiveSwipeIndex int                   `gorm:"column:active_swipe_index;not null;default:0" json:"activeSwipeIndex"`
	SortOrder        int                   `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sortOrder"`
	CreatedAt        time.Time             `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time             `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Swipes []MessageSwipe `gorm:"foreignKey:MessageID" json:"swipes,omitempty"`
}

// TableName 返回表名
func (Message) TableName() string {
	return "messages"
}

// MessageSwipe 消息内容表实体(Swipe)
type MessageSwipe struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MessageID int       `gorm:"column:message_id;not null;index:idx_message_id" json:"messageId"`
	Content   string    `gorm:"column:content;type:text;not null" json:"content"`
	SortOrder int       `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sortOrder"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

// TableName 返回表名
func (MessageSwipe) TableName() string {
	return "message_swipes"
}
