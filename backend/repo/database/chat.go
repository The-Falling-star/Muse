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
func (c *ChatRepo) CreateSession(session *entity.ChatSession) *connect.Error {
	db := config.GetDB()
	result := db.Create(session)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建聊天会话失败: %v", result.Error)
	}
	return nil
}

// GetSessionByID 根据ID获取聊天会话
func (c *ChatRepo) GetSessionByID(id int, userID int) (*entity.ChatSession, *connect.Error) {
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
func (c *ChatRepo) ListSessions(userID int, characterID int, page int, pageSize int) ([]*entity.ChatSession, int64, *connect.Error) {
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
func (c *ChatRepo) UpdateSession(session *entity.ChatSession) *connect.Error {
	db := config.GetDB()

	// 使用乐观锁更新
	result := db.Model(session).
		Where("id = ? AND user_id = ? AND version = ?", session.ID, session.UserID, session.Version).
		Updates(map[string]interface{}{
			"name":    session.Name,
			"version": session.Version + 1,
		})

	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新聊天会话失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeAborted, "更新聊天会话失败：记录不存在或版本号不匹配")
	}

	// 更新内存中的版本号
	session.Version++
	return nil
}

// DeleteSession 删除聊天会话
func (c *ChatRepo) DeleteSession(id int, userID int) *connect.Error {
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
func (c *ChatRepo) GetMessageByID(id int) (*entity.Message, *connect.Error) {
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
func (c *ChatRepo) DeleteMessage(id int) *connect.Error {
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
func (c *ChatRepo) UpdateMessage(message *entity.Message) *connect.Error {
	db := config.GetDB()
	result := db.Model(message).
		Where("id = ?", message.ID).
		Update("active_swipe_index", message.ActiveSwipeIndex)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新消息的active swipe失败: %v", result.Error)
	}
	return nil
}
