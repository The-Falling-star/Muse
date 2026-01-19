package database

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// RegexRuleRepo 正则规则数据库仓库
type RegexRuleRepo struct {
}

// Create 创建正则规则
func (r *RegexRuleRepo) Create(rule *entity.RegexRule) *connect.Error {
	db := config.GetDB()
	result := db.Create(rule)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建正则规则失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据ID获取正则规则
func (r *RegexRuleRepo) GetByID(id int) (*entity.RegexRule, *connect.Error) {
	db := config.GetDB()
	var rule entity.RegexRule
	result := db.Where("id = ?", id).First(&rule)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取正则规则失败: %v", result.Error)
	}
	return &rule, nil
}

// List 获取预设的正则规则列表
func (r *RegexRuleRepo) List(presetID int) ([]*entity.RegexRule, *connect.Error) {
	db := config.GetDB()
	var rules []*entity.RegexRule
	result := db.Where("preset_id = ?", presetID).
		Order("sort_order ASC").
		Find(&rules)
	if result.Error != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取正则规则列表失败: %v", result.Error)
	}
	return rules, nil
}

// Update 更新正则规则
func (r *RegexRuleRepo) Update(rule *entity.RegexRule) *connect.Error {
	db := config.GetDB()
	result := db.Model(rule).
		Where("id = ?", rule.ID).
		Updates(map[string]interface{}{
			"name":                       rule.Name,
			"find_pattern":               rule.FindPattern,
			"replace_pattern":            rule.ReplacePattern,
			"is_enabled":                 rule.IsEnabled,
			"run_on_edit":                rule.RunOnEdit,
			"substitute_regex":           rule.SubstituteRegex,
			"min_depth":                  rule.MinDepth,
			"max_depth":                  rule.MaxDepth,
			"affect_flags_user_input":    rule.AffectFlagsUserInput,
			"affect_flags_ai_output":     rule.AffectFlagsAIOutput,
			"affect_flags_slash_command": rule.AffectFlagsSlashCommand,
			"affect_flags_world_info":    rule.AffectFlagsWorldInfo,
			"affect_flags_prompt":        rule.AffectFlagsPrompt,
			"sort_order":                 rule.SortOrder,
		})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新正则规则失败: %v", result.Error)
	}
	return nil
}

// Delete 删除正则规则
func (r *RegexRuleRepo) Delete(id int) *connect.Error {
	db := config.GetDB()
	result := db.Where("id = ?", id).Delete(&entity.RegexRule{})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除正则规则失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "删除正则规则失败：记录不存在")
	}
	return nil
}

// UpdateRulesOrder 更新正则规则排序
func (r *RegexRuleRepo) UpdateRulesOrder(presetID int, ruleOrders map[int]int) *connect.Error {
	db := config.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for ruleID, sortOrder := range ruleOrders {
		if err := tx.Model(&entity.RegexRule{}).
			Where("id = ? AND preset_id = ?", ruleID, presetID).
			Update("sort_order", sortOrder).Error; err != nil {
			tx.Rollback()
			return errs.NewStandardf(connect.CodeInternal, "更新正则规则排序失败: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新正则规则排序失败：提交事务时出错: %v", err)
	}
	return nil
}

// BatchCreate 批量创建正则规则
func (r *RegexRuleRepo) BatchCreate(rules []*entity.RegexRule) *connect.Error {
	if len(rules) == 0 {
		return nil
	}
	db := config.GetDB()
	result := db.Create(rules)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "批量创建正则规则失败: %v", result.Error)
	}
	return nil
}

// GetMaxSortOrder 获取预设下正则规则的最大排序号
// presetID 为 0 时表示全局正则
func (r *RegexRuleRepo) GetMaxSortOrder(presetID int) (int, *connect.Error) {
	db := config.GetDB()
	var maxOrder int
	result := db.Model(&entity.RegexRule{}).
		Where("preset_id = ?", presetID).
		Select("COALESCE(MAX(sort_order), -1)").
		Scan(&maxOrder)
	if result.Error != nil {
		return 0, errs.NewStandardf(connect.CodeInternal, "获取最大排序号失败: %v", result.Error)
	}
	return maxOrder, nil
}
