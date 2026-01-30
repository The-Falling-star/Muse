package user

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/crypto"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/database"
	"golang.org/x/crypto/bcrypt"
)

type userImpl struct {
	userRepo *database.UserRepo
}

func newUser() *userImpl {
	return &userImpl{
		userRepo: &database.UserRepo{},
	}
}

// GetPublicConfig 获取公共配置（无需认证）
func (u *userImpl) GetPublicConfig(ctx context.Context, req *pb.GetPublicConfigRequest) (*pb.GetPublicConfigResponse, error) {
	return &pb.GetPublicConfigResponse{
		SkipAuth: config.Get().Auth.SkipAuth,
	}, nil
}

func (u *userImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 检查权限
	if userId := jwt.GetUserId(ctx); userId != config.Get().Auth.AdminUserId {
		return nil, errs.NewStandard(connect.CodePermissionDenied, "请使用管理员创建用户")
	}

	// 参数校验
	username := strings.TrimSpace(req.GetUsername())
	if username == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyUsername)
	}
	if req.GetPassword() == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPassword)
	}

	// 检查用户是否已存在
	existingUser, _ := u.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errs.NewStandard(connect.CodeAlreadyExists, errs.UserAlreadyExists)
	}

	// 密码加密
	hashedPassword, err := hashPassword(req.GetPassword())
	if err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "密码加密失败: %v", err)
	}

	// 创建用户
	user := &entity.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 生成JWT Token
	token, err := jwt.GenerateJWT(user.ID)
	if err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "生成Token失败: %v", err)
	}

	return &pb.RegisterResponse{
		User:  convert.UserEntityToPb(user),
		Token: token,
	}, nil
}

// Login 登录
func (u *userImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 参数校验
	username := strings.TrimSpace(req.GetUsername())
	if username == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyUsername)
	}
	if req.GetPassword() == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPassword)
	}

	// 查找用户
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.NewStandard(connect.CodeUnauthenticated, errs.UserNotFound)
	}

	// 验证密码
	if err := verifyPassword(user.PasswordHash, req.GetPassword()); err != nil {
		return nil, errs.NewStandard(connect.CodeUnauthenticated, errs.InvalidPassword)
	}

	// 生成JWT Token
	token, genErr := jwt.GenerateJWT(user.ID)
	if genErr != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "生成Token失败: %v", genErr)
	}
	return &pb.LoginResponse{
		User:  convert.UserEntityToPb(user),
		Token: token,
	}, nil
}

func (u *userImpl) GetCurrentUser(ctx context.Context, req *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 查找用户
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.UserNotFound)
	}

	return &pb.GetCurrentUserResponse{
		User: convert.UserEntityToPb(user),
	}, nil
}

func (u *userImpl) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	if req.GetOldPassword() == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidOldPassword)
	}
	if req.GetNewPassword() == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPassword)
	}

	// 查找用户
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.UserNotFound)
	}

	// 验证旧密码
	if err := verifyPassword(user.PasswordHash, req.GetOldPassword()); err != nil {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidOldPassword)
	}

	// 加密新密码
	hashedPassword, hashErr := hashPassword(req.GetNewPassword())
	if hashErr != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "密码加密失败: %v", hashErr)
	}

	// 更新密码
	user.PasswordHash = hashedPassword
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &pb.ChangePasswordResponse{}, nil
}

func (u *userImpl) ListPersonas(ctx context.Context, req *pb.ListPersonasRequest) (*pb.ListPersonasResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 获取人设列表
	personas, err := u.userRepo.ListPersonas(userID)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbPersonas := make([]*pb.Persona, 0, len(personas))
	for _, persona := range personas {
		pbPersonas = append(pbPersonas, convert.PersonaEntityToPb(persona))
	}

	return &pb.ListPersonasResponse{
		Personas: pbPersonas,
	}, nil
}

func (u *userImpl) GetPersona(ctx context.Context, req *pb.GetPersonaRequest) (*pb.GetPersonaResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPersonaID)
	}

	// 获取人设
	persona, err := u.userRepo.GetPersonaByID(id, userID)
	if err != nil {
		return nil, err
	}
	if persona == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.PersonaNotFound)
	}

	return &pb.GetPersonaResponse{
		Persona: convert.PersonaEntityToPb(persona),
	}, nil
}

func (u *userImpl) CreatePersona(ctx context.Context, req *pb.CreatePersonaRequest) (*pb.CreatePersonaResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPersonaName)
	}

	// 创建人设
	persona := &entity.Persona{
		UserID:      userID,
		Name:        name,
		Avatar:      req.GetAvatar(),
		Description: req.GetDescription(),
	}

	if err := u.userRepo.CreatePersona(persona); err != nil {
		return nil, err
	}

	return &pb.CreatePersonaResponse{
		Persona: convert.PersonaEntityToPb(persona),
	}, nil
}

func (u *userImpl) UpdatePersona(ctx context.Context, req *pb.UpdatePersonaRequest) (*pb.UpdatePersonaResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPersonaID)
	}
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyPersonaName)
	}

	// 获取人设
	persona, err := u.userRepo.GetPersonaByID(id, userID)
	if err != nil {
		return nil, err
	}
	if persona == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.PersonaNotFound)
	}

	// 更新人设字段
	persona.Name = name
	if req.Avatar != nil {
		persona.Avatar = *req.Avatar
	}
	if req.Description != nil {
		persona.Description = *req.Description
	}

	if err = u.userRepo.UpdatePersona(persona); err != nil {
		return nil, err
	}

	return &pb.UpdatePersonaResponse{
		Persona: convert.PersonaEntityToPb(persona),
	}, nil
}

func (u *userImpl) DeletePersona(ctx context.Context, req *pb.DeletePersonaRequest) (*pb.DeletePersonaResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPersonaID)
	}

	// 删除人设
	if err := u.userRepo.DeletePersona(id, userID); err != nil {
		return nil, err
	}

	return &pb.DeletePersonaResponse{}, nil
}

func (u *userImpl) SetActivePersona(ctx context.Context, req *pb.SetActivePersonaRequest) (*pb.SetActivePersonaResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	personaID := int(req.GetPersonaId())
	if personaID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidPersonaID)
	}

	// 验证人设是否存在
	_, err := u.userRepo.GetPersonaByID(personaID, userID)
	if err != nil {
		return nil, err
	}

	// 更新用户的活跃人设ID
	if err := u.userRepo.UpdateActivePersonaID(userID, personaID); err != nil {
		return nil, err
	}

	return &pb.SetActivePersonaResponse{}, nil
}

func (u *userImpl) GetUserSetting(ctx context.Context, req *pb.GetUserSettingRequest) (*pb.GetUserSettingResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 获取用户设置
	setting, err := u.userRepo.GetUserSettingByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 如果用户设置不存在，创建默认设置
	if setting == nil {
		setting = &entity.UserSetting{
			UserID:         userID,
			Theme:          0, // Auto
			Language:       "zh-CN",
			SendOnEnter:    true,
			ShowTimestamps: true,
		}
		if err = u.userRepo.CreateUserSetting(setting); err != nil {
			return nil, err
		}
	}

	return &pb.GetUserSettingResponse{
		Setting: convert.UserSettingEntityToPb(setting),
	}, nil
}

func (u *userImpl) UpdateUserSetting(ctx context.Context, req *pb.UpdateUserSettingRequest) (*pb.UpdateUserSettingResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 获取用户设置
	setting, err := u.userRepo.GetUserSettingByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 如果用户设置不存在，创建新的设置
	if setting == nil {
		setting = &entity.UserSetting{
			UserID: userID,
		}
	}

	// 更新设置字段
	setting.Theme = int(req.GetTheme())
	setting.Language = req.GetLanguage()
	setting.SendOnEnter = req.GetSendOnEnter()
	setting.ShowTimestamps = req.GetShowTimestamps()

	if setting.ID == 0 {
		if err = u.userRepo.CreateUserSetting(setting); err != nil {
			return nil, err
		}
	} else {
		if err = u.userRepo.UpdateUserSetting(setting); err != nil {
			return nil, err
		}
	}

	return &pb.UpdateUserSettingResponse{
		Setting: convert.UserSettingEntityToPb(setting),
	}, nil
}

func (u *userImpl) ListAPIConfigs(ctx context.Context, req *pb.ListAPIConfigsRequest) (*pb.ListAPIConfigsResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 获取API配置列表
	configs, err := u.userRepo.ListAPIConfigs(userID)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式（不包含API Key）
	pbConfigs := make([]*pb.APIConfig, 0, len(configs))
	for _, apiConfig := range configs {
		pbAPIConfig := convert.APIConfigEntityToPb(apiConfig)
		if !config.Get().APIEncrypt.AllowGetKey {
			pbAPIConfig.ApiKey = ""
		}
		pbConfigs = append(pbConfigs, pbAPIConfig)
	}

	return &pb.ListAPIConfigsResponse{
		Configs: pbConfigs,
	}, nil
}

func (u *userImpl) GetAPIConfig(ctx context.Context, req *pb.GetAPIConfigRequest) (*pb.GetAPIConfigResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidAPIConfigID)
	}

	// 获取API配置
	apiConfig, err := u.userRepo.GetAPIConfigByID(id, userID)
	if err != nil {
		return nil, err
	}
	if apiConfig == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.APIConfigNotFound)
	}

	pbAPIConfig := convert.APIConfigEntityToPb(apiConfig)
	if !config.Get().APIEncrypt.AllowGetKey {
		pbAPIConfig.ApiKey = ""
	}
	return &pb.GetAPIConfigResponse{
		Config: pbAPIConfig,
	}, nil
}

func (u *userImpl) CreateAPIConfig(ctx context.Context, req *pb.CreateAPIConfigRequest) (*pb.CreateAPIConfigResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyAPIConfigName)
	}
	if req.GetApiKey() == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyAPIKey)
	}
	if req.GetProvider() == pb.APIProvider_APIProviderUnspecified {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidAPIProvider)
	}

	// 获取API密钥加密配置
	apiEncryptCfg := config.Get().APIEncrypt
	var encryptedAPIKey string

	aesCrypt, err := crypto.NewAESGCMFromKey(apiEncryptCfg.EncryptionKey)
	if err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "初始化加密器失败: %v", err)
	}
	encryptedAPIKey, err = aesCrypt.Encrypt(req.GetApiKey())
	if err != nil {
		return nil, errs.NewStandardf(connect.CodeInternal, "加密API Key失败: %v", err)
	}

	// 创建API配置
	apiConfig := &entity.APIConfig{
		UserID:   userID,
		Name:     name,
		Provider: req.GetProvider(),
		APIKey:   encryptedAPIKey,
		BaseURL:  req.GetBaseUrl(),
		Model:    req.GetModel(),
		IsActive: false, // 默认不激活
	}

	if err = u.userRepo.CreateAPIConfig(apiConfig); err != nil {
		return nil, err
	}

	return &pb.CreateAPIConfigResponse{
		Config: convert.APIConfigEntityToPbWithKey(apiConfig),
	}, nil
}

func (u *userImpl) UpdateAPIConfig(ctx context.Context, req *pb.UpdateAPIConfigRequest) (*pb.UpdateAPIConfigResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidAPIConfigID)
	}
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptyAPIConfigName)
	}

	// 获取API配置
	apiConfig, err := u.userRepo.GetAPIConfigByID(id, userID)
	if err != nil {
		return nil, err
	}
	if apiConfig == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.APIConfigNotFound)
	}

	// 获取API密钥加密配置
	apiEncryptCfg := config.Get().APIEncrypt

	// 更新API配置字段
	apiConfig.Name = name
	apiConfig.Provider = req.GetProvider()
	if req.ApiKey != nil {
		aesCrypt, err := crypto.NewAESGCMFromKey(apiEncryptCfg.EncryptionKey)
		if err != nil {
			return nil, errs.NewStandardf(connect.CodeInternal, "初始化加密器失败: %v", err)
		}
		encryptedAPIKey, err := aesCrypt.Encrypt(*req.ApiKey)
		if err != nil {
			return nil, errs.NewStandardf(connect.CodeInternal, "加密API Key失败: %v", err)
		}
		apiConfig.APIKey = encryptedAPIKey
	}
	if req.BaseUrl != nil {
		apiConfig.BaseURL = *req.BaseUrl
	}
	if req.Model != nil {
		apiConfig.Model = *req.Model
	}

	if err = u.userRepo.UpdateAPIConfig(apiConfig); err != nil {
		return nil, err
	}

	return &pb.UpdateAPIConfigResponse{
		Config: convert.APIConfigEntityToPbWithKey(apiConfig),
	}, nil
}

func (u *userImpl) DeleteAPIConfig(ctx context.Context, req *pb.DeleteAPIConfigRequest) (*pb.DeleteAPIConfigResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidAPIConfigID)
	}

	// 删除API配置
	if err := u.userRepo.DeleteAPIConfig(id, userID); err != nil {
		return nil, err
	}

	return &pb.DeleteAPIConfigResponse{}, nil
}

func (u *userImpl) SetActiveAPIConfig(ctx context.Context, req *pb.SetActiveAPIConfigRequest) (*pb.SetActiveAPIConfigResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	configID := int(req.GetConfigId())
	if configID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidAPIConfigID)
	}

	// 激活API配置
	if err := u.userRepo.ActivateAPIConfig(configID, userID); err != nil {
		return nil, err
	}

	return &pb.SetActiveAPIConfigResponse{}, nil
}

func (u *userImpl) TestAPIConfig(ctx context.Context, req *pb.TestAPIConfigRequest) (*pb.TestAPIConfigResponse, error) {
	// 获取当前用户ID
	userID := jwt.GetUserId(ctx)

	// 参数校验
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidAPIConfigID)
	}

	// 获取API配置
	apiConfig, err := u.userRepo.GetAPIConfigByID(id, userID)
	if err != nil {
		return nil, err
	}
	if apiConfig == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.APIConfigNotFound)
	}

	// TODO: 这里需要实现实际的API测试逻辑
	// 根据不同的provider调用不同的API进行测试
	// 目前返回一个简单的测试结果

	// 模拟测试成功
	return &pb.TestAPIConfigResponse{
		Success:   true,
		ModelInfo: &apiConfig.Model,
	}, nil
}

// hashPassword 对密码进行哈希
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// verifyPassword 验证密码
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
