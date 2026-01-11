package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/user"
)

// UserServer 用户服务（包含认证、人设、用户设置和API配置的管理）
type UserServer struct {
	user user.User
}

// NewUserServer 创建一个新的UserServer实例
func NewUserServer() *UserServer {
	return &UserServer{
		user: user.NewUser(),
	}
}

// Register 处理用户注册请求
func (u *UserServer) Register(ctx context.Context, req *connect.Request[pb.RegisterRequest]) (
	*connect.Response[pb.RegisterResponse], error) {
	resp, err := u.user.Register(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// Login 处理用户登录请求
func (u *UserServer) Login(ctx context.Context, req *connect.Request[pb.LoginRequest]) (
	*connect.Response[pb.LoginResponse], error) {
	resp, err := u.user.Login(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetCurrentUser 获取当前登录用户信息
func (u *UserServer) GetCurrentUser(ctx context.Context, req *connect.Request[pb.GetCurrentUserRequest]) (
	*connect.Response[pb.GetCurrentUserResponse], error) {
	resp, err := u.user.GetCurrentUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ChangePassword 处理用户修改密码请求
func (u *UserServer) ChangePassword(ctx context.Context, req *connect.Request[pb.ChangePasswordRequest]) (
	*connect.Response[pb.ChangePasswordResponse], error) {
	resp, err := u.user.ChangePassword(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListPersonas 获取人设列表
func (u *UserServer) ListPersonas(ctx context.Context, req *connect.Request[pb.ListPersonasRequest]) (
	*connect.Response[pb.ListPersonasResponse], error) {
	resp, err := u.user.ListPersonas(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetPersona 获取指定人设详情
func (u *UserServer) GetPersona(ctx context.Context, req *connect.Request[pb.GetPersonaRequest]) (
	*connect.Response[pb.GetPersonaResponse], error) {
	resp, err := u.user.GetPersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// CreatePersona 创建新人设
func (u *UserServer) CreatePersona(ctx context.Context, req *connect.Request[pb.CreatePersonaRequest]) (
	*connect.Response[pb.CreatePersonaResponse], error) {
	resp, err := u.user.CreatePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdatePersona 更新指定人设
func (u *UserServer) UpdatePersona(ctx context.Context, req *connect.Request[pb.UpdatePersonaRequest]) (
	*connect.Response[pb.UpdatePersonaResponse], error) {
	resp, err := u.user.UpdatePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeletePersona 删除指定人设
func (u *UserServer) DeletePersona(ctx context.Context, req *connect.Request[pb.DeletePersonaRequest]) (
	*connect.Response[pb.DeletePersonaResponse], error) {
	resp, err := u.user.DeletePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// SetActivePersona 设置当前激活的人设
func (u *UserServer) SetActivePersona(ctx context.Context, req *connect.Request[pb.SetActivePersonaRequest]) (
	*connect.Response[pb.SetActivePersonaResponse], error) {
	resp, err := u.user.SetActivePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetUserSetting 获取用户设置
func (u *UserServer) GetUserSetting(ctx context.Context, req *connect.Request[pb.GetUserSettingRequest]) (
	*connect.Response[pb.GetUserSettingResponse], error) {
	resp, err := u.user.GetUserSetting(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateUserSetting 更新用户设置
func (u *UserServer) UpdateUserSetting(ctx context.Context, req *connect.Request[pb.UpdateUserSettingRequest]) (
	*connect.Response[pb.UpdateUserSettingResponse], error) {
	resp, err := u.user.UpdateUserSetting(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListAPIConfigs 获取API配置列表
func (u *UserServer) ListAPIConfigs(ctx context.Context, req *connect.Request[pb.ListAPIConfigsRequest]) (
	*connect.Response[pb.ListAPIConfigsResponse], error) {
	resp, err := u.user.ListAPIConfigs(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetAPIConfig 获取指定的API配置
func (u *UserServer) GetAPIConfig(ctx context.Context, req *connect.Request[pb.GetAPIConfigRequest]) (
	*connect.Response[pb.GetAPIConfigResponse], error) {
	resp, err := u.user.GetAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// CreateAPIConfig 创建新的API配置
func (u *UserServer) CreateAPIConfig(ctx context.Context, req *connect.Request[pb.CreateAPIConfigRequest]) (
	*connect.Response[pb.CreateAPIConfigResponse], error) {
	resp, err := u.user.CreateAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateAPIConfig 更新指定的API配置
func (u *UserServer) UpdateAPIConfig(ctx context.Context, req *connect.Request[pb.UpdateAPIConfigRequest]) (
	*connect.Response[pb.UpdateAPIConfigResponse], error) {
	resp, err := u.user.UpdateAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteAPIConfig 删除指定的API配置
func (u *UserServer) DeleteAPIConfig(ctx context.Context, req *connect.Request[pb.DeleteAPIConfigRequest]) (
	*connect.Response[pb.DeleteAPIConfigResponse], error) {
	resp, err := u.user.DeleteAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// SetActiveAPIConfig 设置当前激活的API配置
func (u *UserServer) SetActiveAPIConfig(ctx context.Context, req *connect.Request[pb.SetActiveAPIConfigRequest]) (
	*connect.Response[pb.SetActiveAPIConfigResponse], error) {
	resp, err := u.user.SetActiveAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// TestAPIConfig 测试API配置的连通性
func (u *UserServer) TestAPIConfig(ctx context.Context, req *connect.Request[pb.TestAPIConfigRequest]) (
	*connect.Response[pb.TestAPIConfigResponse], error) {
	resp, err := u.user.TestAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
