package database

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// WorldInfoRepo 世界书数据库仓库
type WorldInfoRepo struct {
}

// Create 创建世界书
func (w *WorldInfoRepo) Create(worldInfo *entity.WorldInfo) *connect.Error {
	db := config.GetDB()
	result := db.Create(worldInfo)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建世界书失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据ID获取世界书
func (w *WorldInfoRepo) GetByID(id int, userID int) (*entity.WorldInfo, *connect.Error) {
	db := config.GetDB()
	var worldInfo entity.WorldInfo
	result := db.Where("id = ? AND user_id = ?", id, userID).
		First(&worldInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取世界书失败: %v", result.Error)
	}
	return &worldInfo, nil
}

// List 获取世界书列表
func (w *WorldInfoRepo) List(userID int, page int, pageSize int) ([]*entity.WorldInfo, int64, *connect.Error) {
	db := config.GetDB()
	var worldInfos []*entity.WorldInfo
	var total int64

	// 计算总数
	if err := db.Model(&entity.WorldInfo{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, errs.NewStandardf(connect.CodeInternal, "查询世界书总数失败: %v", err)
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
		return nil, 0, errs.NewStandardf(connect.CodeInternal, "查询世界书列表失败: %v", result.Error)
	}

	return worldInfos, total, nil
}

// Update 更新世界书
func (w *WorldInfoRepo) Update(worldInfo *entity.WorldInfo) *connect.Error {
	db := config.GetDB()
	result := db.Model(worldInfo).
		Where("id = ? AND user_id = ?", worldInfo.ID, worldInfo.UserID).
		Updates(map[string]interface{}{
			"name":        worldInfo.Name,
			"description": worldInfo.Description,
		})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新世界书失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "更新世界书失败：记录不存在")
	}
	return nil
}

// Delete 删除世界书
func (w *WorldInfoRepo) Delete(id int, userID int) *connect.Error {
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
		return errs.NewStandardf(connect.CodeInternal, "删除世界书失败：删除关联条目时出错: %v", err)
	}

	// 删除世界书
	result := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.WorldInfo{})
	if result.Error != nil {
		tx.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除世界书失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return errs.NewStandard(connect.CodeNotFound, "删除世界书失败：记录不存在")
	}

	if err := tx.Commit().Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除世界书失败：提交事务时出错: %v", err)
	}
	return nil
}

// CreateEntry 创建世界书条目
func (w *WorldInfoRepo) CreateEntry(entry *entity.WorldInfoEntry) *connect.Error {
	db := config.GetDB()
	result := db.Create(entry)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建世界书条目失败: %v", result.Error)
	}
	return nil
}

// BatchCreateEntries 批量创建世界书条目
func (w *WorldInfoRepo) BatchCreateEntries(entries []*entity.WorldInfoEntry) *connect.Error {
	if len(entries) == 0 {
		return nil
	}
	db := config.GetDB()
	result := db.Create(&entries)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "批量创建世界书条目失败: %v", result.Error)
	}
	return nil
}

// GetEntryByID 根据ID获取世界书条目
func (w *WorldInfoRepo) GetEntryByID(id int) (*entity.WorldInfoEntry, *connect.Error) {
	db := config.GetDB()
	var entry entity.WorldInfoEntry
	result := db.Where("id = ?", id).First(&entry)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取世界书条目失败: %v", result.Error)
	}
	return &entry, nil
}

// ListEntries 获取世界书的条目列表
func (w *WorldInfoRepo) ListEntries(worldInfoID int) ([]*entity.WorldInfoEntry, *connect.Error) {
	db := config.GetDB()
	var entries []*entity.WorldInfoEntry
	result := db.Where("world_info_id = ?", worldInfoID).
		Order("sort_order ASC").
		Find(&entries)
	if result.Error != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取世界书条目列表失败: %v", result.Error)
	}
	return entries, nil
}

// UpdateEntry 更新世界书条目
func (w *WorldInfoRepo) UpdateEntry(entry *entity.WorldInfoEntry) *connect.Error {
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
		return errs.NewStandardf(connect.CodeInternal, "更新世界书条目失败: %v", result.Error)
	}
	return nil
}

// DeleteEntry 删除世界书条目
func (w *WorldInfoRepo) DeleteEntry(id int) *connect.Error {
	db := config.GetDB()
	result := db.Where("id = ?", id).Delete(&entity.WorldInfoEntry{})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除世界书条目失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "删除世界书条目失败：记录不存在")
	}
	return nil
}

// UpdateEntriesOrder 更新世界书条目排序
func (w *WorldInfoRepo) UpdateEntriesOrder(worldInfoID int, entryOrders map[int]int) *connect.Error {
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
			return errs.NewStandardf(connect.CodeInternal, "更新世界书条目排序失败: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新世界书条目排序失败：提交事务时出错: %v", err)
	}
	return nil
}

// ListGlobalWorldInfosWithEntries 获取用户所有全局世界书及其启用的条目
func (w *WorldInfoRepo) ListGlobalWorldInfosWithEntries(userID int) ([]*entity.WorldInfo, *connect.Error) {
	db := config.GetDB()
	var worldInfos []*entity.WorldInfo
	result := db.Where("user_id = ? AND is_global = ?", userID, true).
		Preload("Entries", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_enabled = ?", true).Order("sort_order ASC")
		}).
		Find(&worldInfos)
	if result.Error != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取全局世界书列表失败: %v", result.Error)
	}
	return worldInfos, nil
}

// GetByIDWithEntries 根据ID获取世界书及其所有启用的条目
func (w *WorldInfoRepo) GetByIDWithEntries(id int, userID int) (*entity.WorldInfo, *connect.Error) {
	db := config.GetDB()
	var worldInfo entity.WorldInfo
	result := db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Entries", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_enabled = ?", true).Order("sort_order ASC")
		}).
		First(&worldInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取世界书失败: %v", result.Error)
	}
	return &worldInfo, nil
}

// GetVersions 批量获取世界书的版本号
// 用于缓存版本校验，只查询版本字段以减少数据传输
func (w *WorldInfoRepo) GetVersions(ids []int, userID int) (map[int64]int, *connect.Error) {
	if len(ids) == 0 {
		return make(map[int64]int), nil
	}
	db := config.GetDB()
	type versionResult struct {
		ID      int64
		Version int
	}
	var results []versionResult
	err := db.Model(&entity.WorldInfo{}).
		Where("id IN ? AND user_id = ?", ids, userID).
		Select("id", "version").
		Scan(&results).Error
	if err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取世界书版本失败: %v", err)
	}
	versions := make(map[int64]int)
	for _, r := range results {
		versions[r.ID] = r.Version
	}
	return versions, nil
}

// GetGlobalWorldInfoVersions 获取用户所有全局世界书的版本号
func (w *WorldInfoRepo) GetGlobalWorldInfoVersions(userID int) (map[int64]int, *connect.Error) {
	db := config.GetDB()
	type versionResult struct {
		ID      int64
		Version int
	}
	var results []versionResult
	err := db.Model(&entity.WorldInfo{}).
		Where("user_id = ? AND is_global = ?", userID, true).
		Select("id", "version").
		Scan(&results).Error
	if err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取全局世界书版本失败: %v", err)
	}
	versions := make(map[int64]int)
	for _, r := range results {
		versions[r.ID] = r.Version
	}
	return versions, nil
}
