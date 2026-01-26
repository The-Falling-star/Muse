package database

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	"gorm.io/gorm"
)

// UserRepo 用户数据库仓库
type UserRepo struct{}

// Create 创建用户
func (u *UserRepo) Create(user *entity.User) *connect.Error {
	db := config.GetDB()
	result := db.Create(user)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "数据库创建失败: %v", result.Error)
	}
	return nil
}

// GetByID 根据ID获取用户
func (u *UserRepo) GetByID(id int) (*entity.User, *connect.Error) {
	db := config.GetDB()
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
func (u *UserRepo) GetByUsername(username string) (*entity.User, *connect.Error) {
	db := config.GetDB()
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
func (u *UserRepo) Update(user *entity.User) *connect.Error {
	db := config.GetDB()
	result := db.Save(user)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "数据库更新失败: %v", result.Error)
	}
	return nil
}

// UpdateActivePersonaID 更新用户的活跃人设ID
func (u *UserRepo) UpdateActivePersonaID(userID int, personaID int) *connect.Error {
	db := config.GetDB()
	result := db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("active_persona_id", personaID)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新用户失败: %v", result.Error)
	}
	return nil
}

// UpdateActivePresetID 更新用户的活跃预设ID
func (u *UserRepo) UpdateActivePresetID(userID int, presetID int) *connect.Error {
	db := config.GetDB()
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
func (u *UserRepo) CreatePersona(persona *entity.Persona) *connect.Error {
	db := config.GetDB()
	result := db.Create(persona)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建人设失败: %v", result.Error)
	}
	return nil
}

// GetPersonaByID 根据ID获取人设
func (u *UserRepo) GetPersonaByID(id int, userID int) (*entity.Persona, *connect.Error) {
	db := config.GetDB()
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
func (u *UserRepo) ListPersonas(userID int) ([]*entity.Persona, *connect.Error) {
	db := config.GetDB()
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
func (u *UserRepo) UpdatePersona(persona *entity.Persona) *connect.Error {
	db := config.GetDB()
	result := db.Save(persona)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新人设失败: %v", result.Error)
	}
	return nil
}

// DeletePersona 删除人设
func (u *UserRepo) DeletePersona(id int, userID int) *connect.Error {
	db := config.GetDB()
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
func (u *UserRepo) CreateUserSetting(setting *entity.UserSetting) *connect.Error {
	db := config.GetDB()
	result := db.Create(setting)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建用户设置失败: %v", result.Error)
	}
	return nil
}

// GetUserSettingByUserID 根据用户ID获取用户设置
func (u *UserRepo) GetUserSettingByUserID(userID int) (*entity.UserSetting, *connect.Error) {
	db := config.GetDB()
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
func (u *UserRepo) UpdateUserSetting(setting *entity.UserSetting) *connect.Error {
	db := config.GetDB()
	result := db.Save(setting)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新用户设置失败: %v", result.Error)
	}
	return nil
}

// ==================== APIConfig 相关方法 ====================

// CreateAPIConfig 创建API配置
func (u *UserRepo) CreateAPIConfig(apiConfig *entity.APIConfig) *connect.Error {
	db := config.GetDB()
	result := db.Create(apiConfig)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "创建API配置失败: %v", result.Error)
	}
	return nil
}

// GetAPIConfigByID 根据ID获取API配置
func (u *UserRepo) GetAPIConfigByID(id int, userID int) (*entity.APIConfig, *connect.Error) {
	db := config.GetDB()
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
func (u *UserRepo) ListAPIConfigs(userID int) ([]*entity.APIConfig, *connect.Error) {
	db := config.GetDB()
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
func (u *UserRepo) UpdateAPIConfig(apiConfig *entity.APIConfig) *connect.Error {
	db := config.GetDB()
	result := db.Save(apiConfig)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "更新API配置失败: %v", result.Error)
	}
	return nil
}

// DeleteAPIConfig 删除API配置
func (u *UserRepo) DeleteAPIConfig(id int, userID int) *connect.Error {
	db := config.GetDB()
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
func (u *UserRepo) DeactivateAllAPIConfigs(userID int) *connect.Error {
	db := config.GetDB()
	result := db.Model(&entity.APIConfig{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false)
	if result.Error != nil {
		return errs.NewStandardf(connect.CodeInternal, "停用API配置失败: %v", result.Error)
	}
	return nil
}

// ActivateAPIConfig 激活API配置
func (u *UserRepo) ActivateAPIConfig(id int, userID int) *connect.Error {
	db := config.GetDB()
	// 先停用所有配置
	if err := u.DeactivateAllAPIConfigs(userID); err != nil {
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
func (u *UserRepo) GetActiveAPIConfig(userID int) (*entity.APIConfig, *connect.Error) {
	db := config.GetDB()
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
