package database

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PresetRepo 预设数据库仓库
type PresetRepo struct{}

// NewPresetRepo 创建预设数据仓库实例
func NewPresetRepo() PresetRepository {
	return &PresetRepo{}
}

// Create 创建预设
func (p *PresetRepo) Create(ctx context.Context, preset *entity.Preset) error {
	db := GetDB(ctx)
	result := db.Create(preset)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建预设失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据ID获取预设（包含关联的PromptItems）
func (p *PresetRepo) GetByID(ctx context.Context, id int, userID int) (*entity.Preset, error) {
	db := GetDB(ctx)
	var preset entity.Preset
	result := db.Where("id = ? AND user_id = ?", id, userID).
		Preload("PromptItems", func(db *gorm.DB) *gorm.DB { return db.Order("pre") }).
		First(&preset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取预设失败: %v", result.Error)
	}
	return &preset, nil
}

// List 获取预设列表
func (p *PresetRepo) List(ctx context.Context, userID int, page int, pageSize int) (
	[]*entity.Preset, []int, int64, error) {
	db := GetDB(ctx)
	var presets []*entity.Preset
	var total int64

	// 计算总数
	if err := db.Model(&entity.Preset{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, nil, 0, errs.NewStandardf(connect.CodeInternal, "查询预设总数失败: %v", err)
	}

	// 分页查询 - 不加载关联的PromptItems和RegexRules
	offset := (page - 1) * pageSize
	result := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&presets)
	if result.Error != nil {
		return nil, nil, 0, errs.NewStandardf(connect.CodeInternal, "查询预设列表失败: %v", result.Error)
	}

	// 如果没有预设，直接返回
	if len(presets) == 0 {
		return presets, []int{}, total, nil
	}

	ids := make([]int, 0, len(presets))
	for _, preset := range presets {
		ids = append(ids, preset.ID)
	}

	promptLen := make([]struct {
		PresetID  int
		ItemCount int
	}, 0, len(ids))
	if err := db.Model(&entity.PromptItem{}).
		Select("preset_id, COUNT(*) as item_count").
		Where("preset_id IN (?)", ids).
		Group("preset_id").
		Find(&promptLen).Error; err != nil {
		return nil, nil, 0, errs.NewStandardf(connect.CodeInternal, "查询预设提示项数量失败: %v", err)
	}
	log.Debugf("promptLen: %v", promptLen)
	presetIDToPrompt := make(map[int]int, len(promptLen))
	for _, v := range promptLen {
		presetIDToPrompt[v.PresetID] = v.ItemCount
	}
	presetPromptLen := make([]int, 0, len(presets))
	for _, preset := range presets {
		presetPromptLen = append(presetPromptLen, presetIDToPrompt[preset.ID])
	}

	return presets, presetPromptLen, total, nil
}

// Update 更新预设
func (p *PresetRepo) Update(ctx context.Context, preset *entity.Preset) error {
	db := GetDB(ctx)

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
		return errs.NewStandardf(connect.CodeInternal, "更新预设失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeAborted, "更新预设失败：记录不存在或版本号不匹配")
	}

	// 更新内存中的版本号
	preset.Version++
	return nil
}

// Delete 删除预设
func (p *PresetRepo) Delete(ctx context.Context, id int, userID int) error {
	db := GetDB(ctx)
	// 删除关联的提示项
	if err := db.Where("preset_id = ?", id).Delete(&entity.PromptItem{}).Error; err != nil {
		db.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除预设失败：删除关联的提示项时出错: %v", err)
	}

	// 删除预设
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Preset{})
	if result.Error != nil {
		db.Rollback()
		return errs.NewStandardf(connect.CodeInternal, "删除预设失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		db.Rollback()
		return errs.NewStandard(connect.CodeNotFound, "删除预设失败：记录不存在")
	}
	return nil
}

// CreatePromptItem 创建提示项
func (p *PresetRepo) CreatePromptItem(ctx context.Context, item *entity.PromptItem) error {
	db := GetDB(ctx)
	result := db.Create(item)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建提示项失败: %v", result.Error)
	}
	return nil
}

// BatchCreatePromptItem 批量创建提示项
func (p *PresetRepo) BatchCreatePromptItem(ctx context.Context, items []*entity.PromptItem) error {
	db := GetDB(ctx)
	result := db.Create(items)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建提示项失败: %v", result.Error)
	}
	return nil
}

// GetPromptItemByID 根据ID获取提示项
func (p *PresetRepo) GetPromptItemByID(ctx context.Context, id int) (*entity.PromptItem, error) {
	db := GetDB(ctx)
	var item entity.PromptItem
	result := db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取提示项失败: %v", result.Error)
	}
	return &item, nil
}

// ListPromptItems 获取预设的提示项列表
func (p *PresetRepo) ListPromptItems(ctx context.Context, presetID int) ([]*entity.PromptItem, error) {
	db := GetDB(ctx)
	var items []*entity.PromptItem
	result := db.Where("preset_id = ?", presetID).
		Order("sort_order ASC").
		Find(&items)
	if result.Error != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取提示项列表失败: %v", result.Error)
	}
	return items, nil
}

// UpdatePromptItem 更新提示项
func (p *PresetRepo) UpdatePromptItem(ctx context.Context, item *entity.PromptItem) error {
	db := GetDB(ctx)
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
			"pre":                item.Pre,
			"next":               item.Next,
		})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新提示项失败: %v", result.Error)
	}
	return nil
}

// BatchUpdatePromptItemOrder 批量更新提示项排序
func (p *PresetRepo) BatchUpdatePromptItemOrder(ctx context.Context, items []*entity.PromptItem) error {
	db := GetDB(ctx)
	preCase := ""
	nextCase := ""
	ids := make([]int, 0, len(items))
	for _, item := range items {
		pre := "NULL"
		if item.Pre != nil && *item.Pre != 0 {
			pre = fmt.Sprintf("%d", *item.Pre)
		}
		next := "NULL"
		if item.Next != nil && *item.Next != 0 {
			next = fmt.Sprintf("%d", *item.Next)
		}
		preCase += fmt.Sprintf(`WHEN %d THEN %s `, item.ID, pre)
		nextCase += fmt.Sprintf(`WHEN %d THEN %s `, item.ID, next)
		ids = append(ids, item.ID)
	}

	premptItem := entity.PromptItem{}
	sql := fmt.Sprintf(`
UPDATE %s SET
    %s = CASE id
        %s
        ELSE %s
    END,
    %s = CASE id
        %s
        ELSE %s
    END
WHERE id IN (?);`, premptItem.TableName(),
		"pre", preCase, "pre",
		"`next`", nextCase, "`next`")

	if err := db.Exec(sql, ids).Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新提示项排序失败: %v", err)
	}
	return nil
}

// DeletePromptItem 删除提示项
func (p *PresetRepo) DeletePromptItem(ctx context.Context, id int) error {
	db := GetDB(ctx)
	result := db.Where("id = ?", id).Delete(&entity.PromptItem{})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除提示项失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "删除提示项失败：记录不存在")
	}
	return nil
}

// UpdatePromptItemsOrder 更新提示项排序
func (p *PresetRepo) UpdatePromptItemsOrder(ctx context.Context, presetID, userID, sourceID,
	desID int, operation pb.SortOperation) error {
	db := GetDB(ctx)
	var cnt int64 = 0
	if err := db.Where("preset_id = ? AND user_id = ?", presetID, userID).Count(&cnt).Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新提示项排序时获取预设失败: %v", err)
	}
	if cnt == 0 {
		return errs.NewStandard(connect.CodeNotFound, "预设不存在")
	}
	ids := [2]int{sourceID, desID}
	var prompts []*entity.PromptItem
	if err := db.Select("id", "pre", "`next`").
		Where("preset_id = ? AND (id IN(?) OR pre IN (?) OR `next` IN (?))", presetID, ids, ids, ids).
		Find(&prompts).
		Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新提示项排序时获取提示项失败: %v", err)
	}
	promptMap := make(map[int]*entity.PromptItem, len(prompts))
	for _, prompt := range prompts {
		promptMap[prompt.ID] = prompt
	}
	if _, ok := promptMap[sourceID]; !ok {
		return errs.NewStandardf(connect.CodeNotFound, "原提示项 %d 不存在", sourceID)
	}
	if _, ok := promptMap[desID]; !ok {
		return errs.NewStandardf(connect.CodeNotFound, "目标提示项 %d 不存在", desID)
	}
	if operation == pb.SortOperation_Pre {
		moveToHead(promptMap, sourceID, desID)
	} else {
		moveToHead(promptMap, desID, sourceID)
	}
	return nil
}

func moveToHead(prompts map[int]*entity.PromptItem, sourceID, desID int) {
	if prompts[sourceID].Next != nil {
		prompts[*prompts[sourceID].Next].Pre = prompts[sourceID].Pre
	}
	prompts[*prompts[sourceID].Pre].Next = prompts[sourceID].Next

	prompts[sourceID].Pre = prompts[desID].Pre
	prompts[sourceID].Next = &desID

	if prompts[sourceID].Pre != nil {
		prompts[*prompts[sourceID].Pre].Next = &sourceID
	}
	prompts[desID].Pre = &sourceID
}

// GetVersion 获取预设的版本号
// 用于缓存版本校验，只查询版本字段以减少数据传输
func (p *PresetRepo) GetVersion(ctx context.Context, id int, userID int) (int, error) {
	db := GetDB(ctx)
	var version int
	result := db.Model(&entity.Preset{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("version").
		Scan(&version)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errs.NewStandard(connect.CodeNotFound, "预设不存在")
		}
		return 0, errs.NewStandardf(connect.CodeInternal, "获取预设版本失败: %v", result.Error)
	}
	return version, nil
}

// SetActivePreset 设置用户启用的预设
func (p *PresetRepo) SetActivePreset(ctx context.Context, userID, presetID int) error {
	log.Infof("Setting active preset for user %d to %d", userID, presetID)
	db := GetDB(ctx)
	if err := db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("active_preset_id", presetID).Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "设置用户启用的预设失败: %v", err)
	}
	return nil
}
