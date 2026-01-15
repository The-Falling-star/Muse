package database

import (
	"errors"

	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// PresetRepo 预设数据库仓库
type PresetRepo struct {
}

// Create 创建预设
func (p *PresetRepo) Create(preset *entity.Preset) error {
	db := config.GetDB()
	result := db.Create(preset)
	return result.Error
}

// GetByID 根据ID获取预设
func (p *PresetRepo) GetByID(id int, userID int) (*entity.Preset, error) {
	db := config.GetDB()
	var preset entity.Preset
	result := db.Where("id = ? AND user_id = ?", id, userID).
		First(&preset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &preset, nil
}

// List 获取预设列表
func (p *PresetRepo) List(userID int, page int, pageSize int) ([]*entity.Preset, int64, error) {
	db := config.GetDB()
	var presets []*entity.Preset
	var total int64

	// 计算总数
	if err := db.Model(&entity.Preset{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 - 不加载关联的PromptItems和RegexRules
	offset := (page - 1) * pageSize
	result := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&presets)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return presets, total, nil
}

// Update 更新预设
func (p *PresetRepo) Update(preset *entity.Preset) error {
	db := config.GetDB()

	// 使用乐观锁更新
	result := db.Model(preset).
		Where("id = ? AND user_id = ? AND version = ?", preset.ID, preset.UserID, preset.Version).
		Updates(map[string]interface{}{
			"name":              preset.Name,
			"temperature":       preset.Temperature,
			"top_p":             preset.TopP,
			"top_k":             preset.TopK,
			"max_tokens":        preset.MaxTokens,
			"frequency_penalty": preset.FrequencyPenalty,
			"presence_penalty":  preset.PresencePenalty,
			"version":           preset.Version + 1,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("更新失败：记录不存在或版本号不匹配")
	}

	// 更新内存中的版本号
	preset.Version++
	return nil
}

// Delete 删除预设
func (p *PresetRepo) Delete(id int, userID int) error {
	db := config.GetDB()

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除关联的提示项
	if err := tx.Where("preset_id = ?", id).Delete(&entity.PromptItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除预设
	result := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Preset{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("删除失败：记录不存在")
	}

	return tx.Commit().Error
}

// CreatePromptItem 创建提示项
func (p *PresetRepo) CreatePromptItem(item *entity.PromptItem) error {
	db := config.GetDB()
	result := db.Create(item)
	return result.Error
}

// GetPromptItemByID 根据ID获取提示项
func (p *PresetRepo) GetPromptItemByID(id int) (*entity.PromptItem, error) {
	db := config.GetDB()
	var item entity.PromptItem
	result := db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &item, nil
}

// ListPromptItems 获取预设的提示项列表
func (p *PresetRepo) ListPromptItems(presetID int) ([]*entity.PromptItem, error) {
	db := config.GetDB()
	var items []*entity.PromptItem
	result := db.Where("preset_id = ?", presetID).
		Order("sort_order ASC").
		Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

// UpdatePromptItem 更新提示项
func (p *PresetRepo) UpdatePromptItem(item *entity.PromptItem) error {
	db := config.GetDB()
	result := db.Model(item).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"name":               item.Name,
			"content":            item.Content,
			"role":               item.Role,
			"is_enabled":         item.IsEnabled,
			"injection_position": item.InjectionPosition,
			"injection_depth":    item.InjectionDepth,
			"forbid_overrides":   item.ForbidOverrides,
			"sort_order":         item.SortOrder,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeletePromptItem 删除提示项
func (p *PresetRepo) DeletePromptItem(id int) error {
	db := config.GetDB()
	result := db.Where("id = ?", id).Delete(&entity.PromptItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("删除失败：记录不存在")
	}
	return nil
}

// UpdatePromptItemsOrder 更新提示项排序
func (p *PresetRepo) UpdatePromptItemsOrder(presetID int, itemOrders map[int]int) error {
	db := config.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for itemID, sortOrder := range itemOrders {
		if err := tx.Model(&entity.PromptItem{}).
			Where("id = ? AND preset_id = ?", itemID, presetID).
			Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
