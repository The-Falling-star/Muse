package database

import (
	"errors"

	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// WorldInfoRepo 世界书数据库仓库
type WorldInfoRepo struct {
}

// Create 创建世界书
func (w *WorldInfoRepo) Create(worldInfo *entity.WorldInfo) error {
	db := config.GetDB()
	result := db.Create(worldInfo)
	return result.Error
}

// GetByID 根据ID获取世界书
func (w *WorldInfoRepo) GetByID(id int, userID int) (*entity.WorldInfo, error) {
	db := config.GetDB()
	var worldInfo entity.WorldInfo
	result := db.Where("id = ? AND user_id = ?", id, userID).
		First(&worldInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &worldInfo, nil
}

// List 获取世界书列表
func (w *WorldInfoRepo) List(userID int, page int, pageSize int) ([]*entity.WorldInfo, int64, error) {
	db := config.GetDB()
	var worldInfos []*entity.WorldInfo
	var total int64

	// 计算总数
	if err := db.Model(&entity.WorldInfo{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询 - 只查询列表展示需要的字段，不加载关联的Entries
	offset := (page - 1) * pageSize
	result := db.Model(&entity.WorldInfo{}).
		Select("id", "user_id", "name", "is_global", "created_at", "updated_at").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&worldInfos)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return worldInfos, total, nil
}

// Update 更新世界书
func (w *WorldInfoRepo) Update(worldInfo *entity.WorldInfo) error {
	db := config.GetDB()
	result := db.Model(worldInfo).
		Where("id = ? AND user_id = ?", worldInfo.ID, worldInfo.UserID).
		Updates(map[string]interface{}{
			"name":        worldInfo.Name,
			"description": worldInfo.Description,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新失败：记录不存在")
	}
	return nil
}

// Delete 删除世界书
func (w *WorldInfoRepo) Delete(id int, userID int) error {
	db := config.GetDB()

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除关联的条目
	if err := tx.Where("world_info_id = ?", id).Delete(&entity.WorldInfoEntry{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除世界书
	result := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.WorldInfo{})
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

// CreateEntry 创建世界书条目
func (w *WorldInfoRepo) CreateEntry(entry *entity.WorldInfoEntry) error {
	db := config.GetDB()
	result := db.Create(entry)
	return result.Error
}

// GetEntryByID 根据ID获取世界书条目
func (w *WorldInfoRepo) GetEntryByID(id int) (*entity.WorldInfoEntry, error) {
	db := config.GetDB()
	var entry entity.WorldInfoEntry
	result := db.Where("id = ?", id).First(&entry)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &entry, nil
}

// ListEntries 获取世界书的条目列表
func (w *WorldInfoRepo) ListEntries(worldInfoID int) ([]*entity.WorldInfoEntry, error) {
	db := config.GetDB()
	var entries []*entity.WorldInfoEntry
	result := db.Where("world_info_id = ?", worldInfoID).
		Order("sort_order ASC").
		Find(&entries)
	if result.Error != nil {
		return nil, result.Error
	}
	return entries, nil
}

// UpdateEntry 更新世界书条目
func (w *WorldInfoRepo) UpdateEntry(entry *entity.WorldInfoEntry) error {
	db := config.GetDB()
	result := db.Model(entry).
		Where("id = ?", entry.ID).
		Updates(map[string]interface{}{
			"uid":             entry.UID,
			"keys_list":       entry.KeysList,
			"secondary_keys":  entry.SecondaryKeys,
			"content":         entry.Content,
			"comment":         entry.Comment,
			"is_enabled":      entry.IsEnabled,
			"constant":        entry.Constant,
			"selective":       entry.Selective,
			"insertion_order": entry.InsertionOrder,
			"position":        entry.Position,
			"depth":           entry.Depth,
			"sort_order":      entry.SortOrder,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteEntry 删除世界书条目
func (w *WorldInfoRepo) DeleteEntry(id int) error {
	db := config.GetDB()
	result := db.Where("id = ?", id).Delete(&entity.WorldInfoEntry{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("删除失败：记录不存在")
	}
	return nil
}

// UpdateEntriesOrder 更新世界书条目排序
func (w *WorldInfoRepo) UpdateEntriesOrder(worldInfoID int, entryOrders map[int]int) error {
	db := config.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for entryID, sortOrder := range entryOrders {
		if err := tx.Model(&entity.WorldInfoEntry{}).
			Where("id = ? AND world_info_id = ?", entryID, worldInfoID).
			Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
