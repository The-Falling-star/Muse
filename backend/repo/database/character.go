package database

import (
	"errors"

	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// CharacterRepo 角色卡数据库仓库
type CharacterRepo struct {
}

// Create 创建角色
func (c *CharacterRepo) Create(character *entity.Character) error {
	db := config.GetDB()
	result := db.Create(character)
	return result.Error
}

// GetByID 根据ID获取角色
func (c *CharacterRepo) GetByID(id int, userID int) (*entity.Character, error) {
	db := config.GetDB()
	var character entity.Character
	result := db.Where("id = ? AND user_id = ?", id, userID).
		First(&character)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &character, nil
}

// List 获取角色列表
func (c *CharacterRepo) List(userID int, page int, pageSize int) ([]*entity.Character, int64, error) {
	db := config.GetDB()
	var characters []*entity.Character
	var total int64

	// 计算总数
	if err := db.Model(&entity.Character{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 - 只查询列表展示需要的字段
	offset := (page - 1) * pageSize
	result := db.Model(&entity.Character{}).
		Where("user_id = ?", userID).
		Select("id", "user_id", "name", "avatar", "world_info_id", "lock_version", "created_at", "updated_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&characters)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return characters, total, nil
}

// Update 更新角色
func (c *CharacterRepo) Update(character *entity.Character) error {
	db := config.GetDB()

	// 使用乐观锁更新
	result := db.Model(character).
		Where("id = ? AND user_id = ? AND lock_version = ?", character.ID, character.UserID, character.Version).
		Updates(map[string]interface{}{
			"name":             character.Name,
			"avatar":           character.Avatar,
			"description":      character.Description,
			"first_message":    character.FirstMessage,
			"example_dialogue": character.ExampleDialogue,
			"creator_notes":    character.CreatorNotes,
			"world_info_id":    character.WorldInfoID,
			"lock_version":     character.Version + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("更新失败：记录不存在或版本号不匹配")
	}

	// 更新内存中的版本号
	character.Version++
	return nil
}

// Delete 删除角色
func (c *CharacterRepo) Delete(id int, userID int) error {
	db := config.GetDB()
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Character{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("删除失败：记录不存在")
	}
	return nil
}
