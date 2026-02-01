package database

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// ChatRepo 聊天会话数据库仓库
type ChatRepo struct {
}

// CreateSession 创建聊天会话
func (c *ChatRepo) CreateSession(session *entity.ChatSession) error {
	db := config.GetDB()
	result := db.Create(session)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建聊天会话失败: %v", result.Error)
	}
	return nil
}

// GetSessionByID 根据ID获取聊天会话
func (c *ChatRepo) GetSessionByID(id int, userID int) (*entity.ChatSession, error) {
	db := config.GetDB()
	var session entity.ChatSession
	result := db.Where("id = ? AND user_id = ?", id, userID).
		First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取聊天会话失败: %v", result.Error)
	}
	return &session, nil
}

// ListSessions 获取会话列表
func (c *ChatRepo) ListSessions(userID, characterID, page, pageSize int) ([]*entity.ChatSession, int64, error) {
	db := config.GetDB()
	var sessions []*entity.ChatSession
	var total int64

	// 构建查询条件
	query := db.Model(&entity.ChatSession{}).Where("user_id = ?", userID)
	if characterID > 0 {
		query = query.Where("character_id = ?", characterID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errs.NewStandardf(connect.CodeInternal, "查询会话总数失败: %v", err)
	}

	// 分页查询 - 只查询列表展示需要的字段，不加载关联数据
	offset := (page - 1) * pageSize
	query = query.Select("id", "user_id", "character_id", "name", "version", "created_at", "updated_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize)

	result := query.Find(&sessions)
	if result.Error != nil {
		return nil, 0, errs.NewStandardf(connect.CodeInternal, "查询会话列表失败: %v", result.Error)
	}

	return sessions, total, nil
}

// UpdateSession 更新聊天会话
func (c *ChatRepo) UpdateSession(sessionID, userID, version int, sessionName string) (int, error) {
	db := config.GetDB()

	// 使用乐观锁更新
	result := db.Model(&entity.ChatSession{}).
		Where("id = ? AND user_id = ? AND version = ?", sessionID, userID, version).
		Update("name", sessionName).
		Update("version", version+1)

	if result.Error != nil {
		return 0, errs.NewStandardf(connect.CodeInternal, "更新聊天会话失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return 0, errs.NewStandard(connect.CodeAborted, "更新聊天会话失败：记录不存在或版本号不匹配")
	}

	return version + 1, nil
}

// DeleteSession 删除聊天会话
func (c *ChatRepo) DeleteSession(id int, userID int) error {
	db := config.GetDB()

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除关联的消息和swipes
	var session entity.ChatSession
	if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&session).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewStandard(connect.CodeNotFound, "删除会话失败：会话不存在")
		}
		return errs.NewStandardf(connect.CodeInternal, "删除会话失败：查询会话时出错: %v", err)
	}

	// 获取所有消息ID
	var messageIDs []int
	if err := tx.Model(&entity.Message{}).Where("session_id = ?", id).Pluck("id", &messageIDs).Error; err != nil {
		tx.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除会话失败：查询关联消息ID时出错: %v", err)
	}

	// 删除swipes
	if len(messageIDs) > 0 {
		if err := tx.Where("message_id IN ?", messageIDs).Delete(&entity.MessageSwipe{}).Error; err != nil {
			tx.Rollback()
			return errs.NewStandardf(connect.CodeInternal, "删除会话失败：删除关联的swipes时出错: %v", err)
		}
	}

	// 删除消息
	if err := tx.Where("session_id = ?", id).Delete(&entity.Message{}).Error; err != nil {
		tx.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除会话失败：删除关联的消息时出错: %v", err)
	}

	// 删除会话
	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除会话失败：删除会话本体时出错: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除会话失败：提交事务时出错: %v", err)
	}
	return nil
}

// GetMessageByID 根据ID获取消息
func (c *ChatRepo) GetMessageByID(id int) (*entity.Message, error) {
	db := config.GetDB()
	var message entity.Message
	result := db.Where("id = ?", id).
		Preload("Swipes").
		First(&message)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取消息失败: %v", result.Error)
	}
	return &message, nil
}

// DeleteMessage 删除消息
func (c *ChatRepo) DeleteMessage(id int) error {
	db := config.GetDB()

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除swipes
	if err := tx.Where("message_id = ?", id).Delete(&entity.MessageSwipe{}).Error; err != nil {
		tx.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除消息失败：删除关联的swipes时出错: %v", err)
	}

	// 删除消息
	if err := tx.Where("id = ?", id).Delete(&entity.Message{}).Error; err != nil {
		tx.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除消息失败：删除消息本体时出错: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除消息失败：提交事务时出错: %v", err)
	}
	return nil
}

// UpdateMessage 更新消息的activeSwipeIndex
func (c *ChatRepo) UpdateMessage(message *entity.Message) error {
	db := config.GetDB()
	result := db.Model(message).
		Where("id = ?", message.ID).
		Update("active_swipe_index", message.ActiveSwipeIndex)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新消息的active swipe失败: %v", result.Error)
	}
	return nil
}

// GetSessionWithMessages 获取会话及其所有消息和角色卡
func (c *ChatRepo) GetSessionWithMessages(id int, userID int) (*entity.ChatSession, error) {
	db := config.GetDB()
	var session entity.ChatSession
	result := db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Character").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Messages.Swipes", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取聊天会话失败: %v", result.Error)
	}
	return &session, nil
}

// CreateMessage 创建消息
func (c *ChatRepo) CreateMessage(message *entity.Message) error {
	db := config.GetDB()
	result := db.Create(message)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建消息失败: %v", result.Error)
	}
	return nil
}

// CreateMessageSwipe 创建消息swipe
func (c *ChatRepo) CreateMessageSwipe(swipe *entity.MessageSwipe) error {
	db := config.GetDB()
	result := db.Create(swipe)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建消息swipe失败: %v", result.Error)
	}
	return nil
}

// GetMaxMessageSortOrder 获取会话中消息的最大排序号
func (c *ChatRepo) GetMaxMessageSortOrder(sessionID int) (int, error) {
	db := config.GetDB()
	var maxOrder int
	result := db.Model(&entity.Message{}).
		Where("session_id = ?", sessionID).
		Select("COALESCE(MAX(sort_order), -1)").
		Scan(&maxOrder)
	if result.Error != nil {
		return 0, errs.NewStandardf(connect.CodeInternal, "获取最大排序号失败: %v", result.Error)
	}
	return maxOrder, nil
}

// UpdateMessageSwipe 更新消息swipe内容
func (c *ChatRepo) UpdateMessageSwipe(swipeID int, content string) error {
	db := config.GetDB()
	result := db.Model(&entity.MessageSwipe{}).
		Where("id = ?", swipeID).
		Update("content", content)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新消息swipe失败: %v", result.Error)
	}
	return nil
}
