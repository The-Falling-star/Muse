package database

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// CharacterRepo 角色卡数据库仓库
type CharacterRepo struct {
}

// Create 创建角色
func (c *CharacterRepo) Create(character *entity.Character) *connect.Error {
	db := config.GetDB()
	result := db.Create(character)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "数据库创建失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据ID获取角色
func (c *CharacterRepo) GetByID(id int, userID int) (*entity.Character, *connect.Error) {
	db := config.GetDB()
	var character entity.Character
	result := db.Where("id = ? AND user_id = ?", id, userID).
		First(&character)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "数据库查询失败: %v", result.Error)
	}
	return &character, nil
}

// List 获取角色列表
func (c *CharacterRepo) List(userID int, page int, pageSize int) ([]*entity.Character, int64, *connect.Error) {
	db := config.GetDB()
	var characters []*entity.Character
	var total int64

	// 计算总数
	if err := db.Model(&entity.Character{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, errs.NewStandardf(connect.CodeInternal, "数据库查询失败: %v", err)
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
		return nil, 0, errs.NewStandardf(connect.CodeInternal, "数据库查询失败: %v", result.Error)
	}

	return characters, total, nil
}

// Update 更新角色
func (c *CharacterRepo) Update(character *entity.Character) *connect.Error {
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
		return errs.NewStandardf(connect.CodeInternal, "数据库更新失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "角色不存在或已被修改")
	}

	// 更新内存中的版本号
	character.Version++
	return nil
}

// Delete 删除角色
func (c *CharacterRepo) Delete(id int, userID int) *connect.Error {
	db := config.GetDB()
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Character{})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "数据库删除失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "角色不存在")
	}
	return nil
}

// GetVersion 获取角色的版本号
// 用于缓存版本校验，只查询版本字段以减少数据传输
func (c *CharacterRepo) GetVersion(id int, userID int) (int, *connect.Error) {
	db := config.GetDB()
	var version int
	result := db.Model(&entity.Character{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("lock_version").
		Scan(&version)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errs.NewStandard(connect.CodeNotFound, "角色不存在")
		}
		return 0, errs.NewStandardf(connect.CodeInternal, "数据库查询失败: %v", result.Error)
	}
	return version, nil
}
