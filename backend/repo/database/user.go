package database

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// UserRepo 用户数据库仓库
type UserRepo struct{}

// NewUserRepo 创建用户数据仓库实例
func NewUserRepo() UserRepository {
	return &UserRepo{}
}

// Create 创建用户
func (u *UserRepo) Create(ctx context.Context, user *entity.User) error {
	db := GetDB(ctx)
	result := db.Create(user)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "数据库创建失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据ID获取用户
func (u *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	db := GetDB(ctx)
	var user entity.User
	result := db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "数据库查询失败: %v", result.Error)
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (u *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	db := GetDB(ctx)
	var user entity.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "数据库查询失败: %v", result.Error)
	}
	return &user, nil
}

// Update 更新用户
func (u *UserRepo) Update(ctx context.Context, user *entity.User) error {
	db := GetDB(ctx)
	result := db.Save(user)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "数据库更新失败: %v", result.Error)
	}
	return nil
}

// UpdateActivePersonaID 更新用户的活跃人设ID
func (u *UserRepo) UpdateActivePersonaID(ctx context.Context, userID int, personaID int) error {
	db := GetDB(ctx)
	result := db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("active_persona_id", personaID)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新用户失败: %v", result.Error)
	}
	return nil
}

// UpdateActivePresetID 更新用户的活跃预设ID
func (u *UserRepo) UpdateActivePresetID(ctx context.Context, userID int, presetID int) error {
	db := GetDB(ctx)
	result := db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("active_preset_id", presetID)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新用户失败: %v", result.Error)
	}
	return nil
}

// ==================== Persona 相关方法 ====================

// CreatePersona 创建人设
func (u *UserRepo) CreatePersona(ctx context.Context, persona *entity.Persona) error {
	db := GetDB(ctx)
	result := db.Create(persona)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建人设失败: %v", result.Error)
	}
	return nil
}

// GetPersonaByID 根据ID获取人设
func (u *UserRepo) GetPersonaByID(ctx context.Context, id int, userID int) (*entity.Persona, error) {
	db := GetDB(ctx)
	var persona entity.Persona
	result := db.Where("id = ? AND user_id = ?", id, userID).First(&persona)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取人设失败: %v", result.Error)
	}
	return &persona, nil
}

// ListPersonas 获取用户的人设列表
func (u *UserRepo) ListPersonas(ctx context.Context, userID int) ([]*entity.Persona, error) {
	db := GetDB(ctx)
	var personas []*entity.Persona
	result := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&personas)
	if result.Error != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取人设列表失败: %v", result.Error)
	}
	return personas, nil
}

// UpdatePersona 更新人设
func (u *UserRepo) UpdatePersona(ctx context.Context, persona *entity.Persona) error {
	db := GetDB(ctx)
	result := db.Save(persona)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新人设失败: %v", result.Error)
	}
	return nil
}

// DeletePersona 删除人设
func (u *UserRepo) DeletePersona(ctx context.Context, id int, userID int) error {
	db := GetDB(ctx)
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Persona{})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除人设失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "删除人设失败：记录不存在")
	}
	return nil
}

// ==================== UserSetting 相关方法 ====================

// CreateUserSetting 创建用户设置
func (u *UserRepo) CreateUserSetting(ctx context.Context, setting *entity.UserSetting) error {
	db := GetDB(ctx)
	result := db.Create(setting)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建用户设置失败: %v", result.Error)
	}
	return nil
}

// GetUserSettingByUserID 根据用户ID获取用户设置
func (u *UserRepo) GetUserSettingByUserID(ctx context.Context, userID int) (*entity.UserSetting, error) {
	db := GetDB(ctx)
	var setting entity.UserSetting
	result := db.Where("user_id = ?", userID).First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取用户设置失败: %v", result.Error)
	}
	return &setting, nil
}

// UpdateUserSetting 更新用户设置
func (u *UserRepo) UpdateUserSetting(ctx context.Context, setting *entity.UserSetting) error {
	db := GetDB(ctx)
	result := db.Save(setting)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新用户设置失败: %v", result.Error)
	}
	return nil
}

// ==================== APIConfig 相关方法 ====================

// CreateAPIConfig 创建API配置
func (u *UserRepo) CreateAPIConfig(ctx context.Context, apiConfig *entity.APIConfig) error {
	db := GetDB(ctx)
	result := db.Create(apiConfig)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建API配置失败: %v", result.Error)
	}
	return nil
}

// GetAPIConfigByID 根据ID获取API配置
func (u *UserRepo) GetAPIConfigByID(ctx context.Context, id int, userID int) (*entity.APIConfig, error) {
	db := GetDB(ctx)
	var apiConfig entity.APIConfig
	result := db.Where("id = ? AND user_id = ?", id, userID).First(&apiConfig)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取API配置失败: %v", result.Error)
	}
	return &apiConfig, nil
}

// ListAPIConfigs 获取用户的API配置列表
func (u *UserRepo) ListAPIConfigs(ctx context.Context, userID int) ([]*entity.APIConfig, error) {
	db := GetDB(ctx)
	var configs []*entity.APIConfig
	result := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&configs)
	if result.Error != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "获取API配置列表失败: %v", result.Error)
	}
	return configs, nil
}

// UpdateAPIConfig 更新API配置
func (u *UserRepo) UpdateAPIConfig(ctx context.Context, apiConfig *entity.APIConfig) error {
	db := GetDB(ctx)
	result := db.Save(apiConfig)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新API配置失败: %v", result.Error)
	}
	return nil
}

// DeleteAPIConfig 删除API配置
func (u *UserRepo) DeleteAPIConfig(ctx context.Context, id int, userID int) error {
	db := GetDB(ctx)
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.APIConfig{})
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "删除API配置失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "删除API配置失败：记录不存在")
	}
	return nil
}

// DeactivateAllAPIConfigs 停用用户的所有API配置
func (u *UserRepo) DeactivateAllAPIConfigs(ctx context.Context, userID int) error {
	db := GetDB(ctx)
	result := db.Model(&entity.APIConfig{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "停用API配置失败: %v", result.Error)
	}
	return nil
}

// ActivateAPIConfig 激活API配置
func (u *UserRepo) ActivateAPIConfig(ctx context.Context, id int, userID int) error {
	db := GetDB(ctx)
	// 先停用所有配置
	if err := u.DeactivateAllAPIConfigs(ctx, userID); err != nil {
		return err
	}
	// 激活指定配置
	result := db.Model(&entity.APIConfig{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_active", true)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "激活API配置失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errs.NewStandard(connect.CodeNotFound, "激活API配置失败：记录不存在")
	}
	return nil
}

// GetActiveAPIConfig 获取用户当前活跃的API配置
func (u *UserRepo) GetActiveAPIConfig(ctx context.Context, userID int) (*entity.APIConfig, error) {
	db := GetDB(ctx)
	var apiConfig entity.APIConfig
	result := db.Where("user_id = ? AND is_active = ?", userID, true).First(&apiConfig)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.NewStandardf(connect.CodeInternal, "获取活跃API配置失败: %v", result.Error)
	}
	return &apiConfig, nil
}
