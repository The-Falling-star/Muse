package entity

import (
	"time"

	pb "github.com/ling/muse/gen/muse"
)

// ChatSession 聊天会话表实体
// 用于存储用户与角色之间的对话会话
type ChatSession struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                            // 会话唯一标识
	UserID      int       `gorm:"column:user_id;not null;index:idx_user_id" json:"userId"`                 // 所属用户ID
	CharacterID int       `gorm:"column:character_id;not null;index:idx_character_id" json:"characterId"`  // 关联的角色卡ID
	Name        string    `gorm:"column:name;type:varchar(255);not null;default:''" json:"name,omitempty"` // 会话名称，可选，用于用户自定义会话标题
	Version     int       `gorm:"column:version;not null;default:1" json:"version"`                        // 乐观锁版本号，用于并发控制
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                       // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                       // 更新时间

	// 关联关系
	Character *Character `gorm:"foreignKey:CharacterID" json:"character,omitempty"` // 关联的角色卡
	Messages  []Message  `gorm:"foreignKey:SessionID" json:"messages,omitempty"`    // 会话中的消息列表
}

// TableName 返回表名
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// Message 消息表实体
// 表示会话中的一条消息，每条消息可以有多个swipe（候选回复）
type Message struct {
	ID               int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                               // 消息唯一标识
	SessionID        int       `gorm:"column:session_id;not null;index:idx_session_id" json:"sessionId"`           // 所属会话ID
	Role             pb.Role   `gorm:"column:role;type:tinyint unsigned;not null" json:"role"`                     // 消息角色：System/User/Assistant
	ActiveSwipeIndex int       `gorm:"column:active_swipe_index;not null;default:0" json:"activeSwipeIndex"`       // 当前激活的swipe索引，用于在多个候选回复中选择显示哪个
	SortOrder        int       `gorm:"column:sort_order;not null;default:0;index:idx_sort_order" json:"sortOrder"` // 排序顺序，决定消息在会话中的显示顺序
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                          // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                          // 更新时间

	// 关联关系
	Swipes []MessageSwipe `gorm:"foreignKey:MessageID" json:"swipes,omitempty"` // 消息的候选内容列表（swipes）
}

// TableName 返回表名
func (Message) TableName() string {
	return "messages"
}

// MessageSwipe 消息内容表实体(Swipe)
// Swipe是SillyTavern中的概念，表示同一条消息的不同候选回复
// 用户可以在多个候选回复之间切换，类似于"左右滑动"选择不同回复
type MessageSwipe struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                     // Swipe唯一标识
	MessageID int       `gorm:"column:message_id;not null;index:idx_message_id" json:"messageId"` // 所属消息ID
	Content   string    `gorm:"column:content;type:mediumtext;not null;" json:"content"`          // 消息文本内容
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                // 创建时间
}

// TableName 返回表名
func (MessageSwipe) TableName() string {
	return "message_swipes"
}
