package database

import (
	"errors"

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
	return result.Error
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
		return nil, result.Error
	}
	return &session, nil
}

// ListSessions 获取会话列表
func (c *ChatRepo) ListSessions(userID int, characterID int, page int, pageSize int) ([]*entity.ChatSession, int64, error) {
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
		return nil, 0, err
	}

	// 分页查询 - 只查询列表展示需要的字段，不加载关联数据
	offset := (page - 1) * pageSize
	query = query.Select("id", "user_id", "character_id", "name", "version", "created_at", "updated_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize)

	result := query.Find(&sessions)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return sessions, total, nil
}

// UpdateSession 更新聊天会话
func (c *ChatRepo) UpdateSession(session *entity.ChatSession) error {
	db := config.GetDB()

	// 使用乐观锁更新
	result := db.Model(session).
		Where("id = ? AND user_id = ? AND version = ?", session.ID, session.UserID, session.Version).
		Updates(map[string]interface{}{
			"name":    session.Name,
			"version": session.Version + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("更新失败：记录不存在或版本号不匹配")
	}

	// 更新内存中的版本号
	session.Version++
	return nil
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
			return errors.New("会话不存在")
		}
		return err
	}

	// 获取所有消息ID
	var messageIDs []int
	if err := tx.Model(&entity.Message{}).Where("session_id = ?", id).Pluck("id", &messageIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除swipes
	if len(messageIDs) > 0 {
		if err := tx.Where("message_id IN ?", messageIDs).Delete(&entity.MessageSwipe{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 删除消息
	if err := tx.Where("session_id = ?", id).Delete(&entity.Message{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除会话
	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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
		return nil, result.Error
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
		return err
	}

	// 删除消息
	if err := tx.Where("id = ?", id).Delete(&entity.Message{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateMessage 更新消息的activeSwipeIndex
func (c *ChatRepo) UpdateMessage(message *entity.Message) error {
	db := config.GetDB()
	result := db.Model(message).
		Where("id = ?", message.ID).
		Update("active_swipe_index", message.ActiveSwipeIndex)
	return result.Error
}
