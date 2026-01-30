package user

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// User 定义了用户服务的接口（包含认证、人设、用户设置和API配置的管理）
type User interface {
	// GetPublicConfig 获取公共配置（无需认证）
	GetPublicConfig(ctx context.Context, req *pb.GetPublicConfigRequest) (*pb.GetPublicConfigResponse, error)
	// Register 处理用户注册请求
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	// Login 处理用户登录请求
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	// GetCurrentUser 获取当前登录用户信息
	GetCurrentUser(ctx context.Context, req *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error)
	// ChangePassword 处理用户修改密码请求
	ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)

	// ListPersonas 获取人设列表
	ListPersonas(ctx context.Context, req *pb.ListPersonasRequest) (*pb.ListPersonasResponse, error)
	// GetPersona 获取指定人设详情
	GetPersona(ctx context.Context, req *pb.GetPersonaRequest) (*pb.GetPersonaResponse, error)
	// CreatePersona 创建新人设
	CreatePersona(ctx context.Context, req *pb.CreatePersonaRequest) (*pb.CreatePersonaResponse, error)
	// UpdatePersona 更新指定人设
	UpdatePersona(ctx context.Context, req *pb.UpdatePersonaRequest) (*pb.UpdatePersonaResponse, error)
	// DeletePersona 删除指定人设
	DeletePersona(ctx context.Context, req *pb.DeletePersonaRequest) (*pb.DeletePersonaResponse, error)
	// SetActivePersona 设置当前激活的人设
	SetActivePersona(ctx context.Context, req *pb.SetActivePersonaRequest) (*pb.SetActivePersonaResponse, error)

	// GetUserSetting 获取用户设置
	GetUserSetting(ctx context.Context, req *pb.GetUserSettingRequest) (*pb.GetUserSettingResponse, error)
	// UpdateUserSetting 更新用户设置
	UpdateUserSetting(ctx context.Context, req *pb.UpdateUserSettingRequest) (*pb.UpdateUserSettingResponse, error)

	// ListAPIConfigs 获取API配置列表
	ListAPIConfigs(ctx context.Context, req *pb.ListAPIConfigsRequest) (*pb.ListAPIConfigsResponse, error)
	// GetAPIConfig 获取指定的API配置
	GetAPIConfig(ctx context.Context, req *pb.GetAPIConfigRequest) (*pb.GetAPIConfigResponse, error)
	// CreateAPIConfig 创建新的API配置
	CreateAPIConfig(ctx context.Context, req *pb.CreateAPIConfigRequest) (*pb.CreateAPIConfigResponse, error)
	// UpdateAPIConfig 更新指定的API配置
	UpdateAPIConfig(ctx context.Context, req *pb.UpdateAPIConfigRequest) (*pb.UpdateAPIConfigResponse, error)
	// DeleteAPIConfig 删除指定的API配置
	DeleteAPIConfig(ctx context.Context, req *pb.DeleteAPIConfigRequest) (*pb.DeleteAPIConfigResponse, error)
	// SetActiveAPIConfig 设置当前激活的API配置
	SetActiveAPIConfig(ctx context.Context, req *pb.SetActiveAPIConfigRequest) (*pb.SetActiveAPIConfigResponse, error)
	// TestAPIConfig 测试API配置的连通性
	TestAPIConfig(ctx context.Context, req *pb.TestAPIConfigRequest) (*pb.TestAPIConfigResponse, error)
}

// NewUser 创建一个新的User实例
func NewUser() User {
	return newUser()
}
