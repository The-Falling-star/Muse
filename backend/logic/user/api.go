package user

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// User 定义了用户服务的接口（包含认证、人设、用户设置和API配置的管理）
type User interface {
	// Register 处理用户注册请求
	Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRsp, error)
	// Login 处理用户登录请求
	Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRsp, error)
	// GetCurrentUser 获取当前登录用户信息
	GetCurrentUser(ctx context.Context, req *pb.GetCurrentUserReq) (*pb.GetCurrentUserRsp, error)
	// ChangePassword 处理用户修改密码请求
	ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*pb.ChangePasswordRsp, error)

	// ListPersonas 获取人设列表
	ListPersonas(ctx context.Context, req *pb.ListPersonasReq) (*pb.ListPersonasRsp, error)
	// GetPersona 获取指定人设详情
	GetPersona(ctx context.Context, req *pb.GetPersonaReq) (*pb.GetPersonaRsp, error)
	// CreatePersona 创建新人设
	CreatePersona(ctx context.Context, req *pb.CreatePersonaReq) (*pb.CreatePersonaRsp, error)
	// UpdatePersona 更新指定人设
	UpdatePersona(ctx context.Context, req *pb.UpdatePersonaReq) (*pb.UpdatePersonaRsp, error)
	// DeletePersona 删除指定人设
	DeletePersona(ctx context.Context, req *pb.DeletePersonaReq) (*pb.DeletePersonaRsp, error)
	// SetActivePersona 设置当前激活的人设
	SetActivePersona(ctx context.Context, req *pb.SetActivePersonaReq) (*pb.SetActivePersonaRsp, error)

	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, req *pb.GetUserInfoReq) (*pb.GetUserInfoRsp, error)
	// UpdateUserInfo 更新用户信息
	UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoReq) (*pb.UpdateUserInfoRsp, error)

	// ListAPIConfigs 获取API配置列表
	ListAPIConfigs(ctx context.Context, req *pb.ListAPIConfigsReq) (*pb.ListAPIConfigsRsp, error)
	// CreateAPIConfig 创建新的API配置
	CreateAPIConfig(ctx context.Context, req *pb.CreateAPIConfigReq) (*pb.CreateAPIConfigRsp, error)
	// UpdateAPIConfig 更新指定的API配置
	UpdateAPIConfig(ctx context.Context, req *pb.UpdateAPIConfigReq) (*pb.UpdateAPIConfigRsp, error)
	// DeleteAPIConfig 删除指定的API配置
	DeleteAPIConfig(ctx context.Context, req *pb.DeleteAPIConfigReq) (*pb.DeleteAPIConfigRsp, error)
	// SetActiveAPIConfig 设置当前激活的API配置
	SetActiveAPIConfig(ctx context.Context, req *pb.SetActiveAPIConfigReq) (*pb.SetActiveAPIConfigRsp, error)
	// TestAPIConfig 测试API配置的连通性
	TestAPIConfig(ctx context.Context, req *pb.TestAPIConfigReq) (*pb.TestAPIConfigRsp, error)
}

// NewUser 创建一个新的User实例
func NewUser() User {
	return newUser()
}
