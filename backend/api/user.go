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
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "Register", req.Msg, (*pb.RegisterResponse)(nil), err)
	}

	resp, err := u.user.Register(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "Register", req.Msg, resp, err)
	}
	return doResponse(ctx, "Register", req.Msg, resp)
}

// Login 处理用户登录请求
func (u *UserServer) Login(ctx context.Context, req *connect.Request[pb.LoginRequest]) (
	*connect.Response[pb.LoginResponse], error) {
	resp, err := u.user.Login(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "Login", req.Msg, resp, err)
	}
	return doResponse(ctx, "Login", req.Msg, resp)
}

// GetCurrentUser 获取当前登录用户信息
func (u *UserServer) GetCurrentUser(ctx context.Context, req *connect.Request[pb.GetCurrentUserRequest]) (
	*connect.Response[pb.GetCurrentUserResponse], error) {
	resp, err := u.user.GetCurrentUser(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetCurrentUser", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetCurrentUser", req.Msg, resp)
}

// ChangePassword 处理用户修改密码请求
func (u *UserServer) ChangePassword(ctx context.Context, req *connect.Request[pb.ChangePasswordRequest]) (
	*connect.Response[pb.ChangePasswordResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "ChangePassword", req.Msg, (*pb.ChangePasswordResponse)(nil), err)
	}

	resp, err := u.user.ChangePassword(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ChangePassword", req.Msg, resp, err)
	}
	return doResponse(ctx, "ChangePassword", req.Msg, resp)
}

// ListPersonas 获取人设列表
func (u *UserServer) ListPersonas(ctx context.Context, req *connect.Request[pb.ListPersonasRequest]) (
	*connect.Response[pb.ListPersonasResponse], error) {
	resp, err := u.user.ListPersonas(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ListPersonas", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListPersonas", req.Msg, resp)
}

// GetPersona 获取指定人设详情
func (u *UserServer) GetPersona(ctx context.Context, req *connect.Request[pb.GetPersonaRequest]) (
	*connect.Response[pb.GetPersonaResponse], error) {
	resp, err := u.user.GetPersona(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetPersona", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetPersona", req.Msg, resp)
}

// CreatePersona 创建新人设
func (u *UserServer) CreatePersona(ctx context.Context, req *connect.Request[pb.CreatePersonaRequest]) (
	*connect.Response[pb.CreatePersonaResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "CreatePersona", req.Msg, (*pb.CreatePersonaResponse)(nil), err)
	}

	resp, err := u.user.CreatePersona(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "CreatePersona", req.Msg, resp, err)
	}
	return doResponse(ctx, "CreatePersona", req.Msg, resp)
}

// UpdatePersona 更新指定人设
func (u *UserServer) UpdatePersona(ctx context.Context, req *connect.Request[pb.UpdatePersonaRequest]) (
	*connect.Response[pb.UpdatePersonaResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdatePersona", req.Msg, (*pb.UpdatePersonaResponse)(nil), err)
	}

	resp, err := u.user.UpdatePersona(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdatePersona", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdatePersona", req.Msg, resp)
}

// DeletePersona 删除指定人设
func (u *UserServer) DeletePersona(ctx context.Context, req *connect.Request[pb.DeletePersonaRequest]) (
	*connect.Response[pb.DeletePersonaResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "DeletePersona", req.Msg, (*pb.DeletePersonaResponse)(nil), err)
	}

	resp, err := u.user.DeletePersona(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeletePersona", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeletePersona", req.Msg, resp)
}

// SetActivePersona 设置当前激活的人设
func (u *UserServer) SetActivePersona(ctx context.Context, req *connect.Request[pb.SetActivePersonaRequest]) (
	*connect.Response[pb.SetActivePersonaResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "SetActivePersona", req.Msg, (*pb.SetActivePersonaResponse)(nil), err)
	}

	resp, err := u.user.SetActivePersona(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "SetActivePersona", req.Msg, resp, err)
	}
	return doResponse(ctx, "SetActivePersona", req.Msg, resp)
}

// GetUserInfo 获取用户信息
func (u *UserServer) GetUserInfo(ctx context.Context, req *connect.Request[pb.GetUserInfoRequest]) (
	*connect.Response[pb.GetUserInfoResponse], error) {
	resp, err := u.user.GetUserInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetUserInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetUserInfo", req.Msg, resp)
}

// UpdateUserInfo 更新用户信息
func (u *UserServer) UpdateUserInfo(ctx context.Context, req *connect.Request[pb.UpdateUserInfoRequest]) (
	*connect.Response[pb.UpdateUserInfoResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdateUserInfo", req.Msg, (*pb.UpdateUserInfoResponse)(nil), err)
	}

	resp, err := u.user.UpdateUserInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateUserInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateUserInfo", req.Msg, resp)
}

// ListAPIConfigs 获取API配置列表
func (u *UserServer) ListAPIConfigs(ctx context.Context, req *connect.Request[pb.ListAPIConfigsRequest]) (
	*connect.Response[pb.ListAPIConfigsResponse], error) {
	resp, err := u.user.ListAPIConfigs(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ListAPIConfigs", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListAPIConfigs", req.Msg, resp)
}

// CreateAPIConfig 创建新的API配置
func (u *UserServer) CreateAPIConfig(ctx context.Context, req *connect.Request[pb.CreateAPIConfigRequest]) (
	*connect.Response[pb.CreateAPIConfigResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "CreateAPIConfig", req.Msg, (*pb.CreateAPIConfigResponse)(nil), err)
	}

	resp, err := u.user.CreateAPIConfig(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "CreateAPIConfig", req.Msg, resp, err)
	}
	return doResponse(ctx, "CreateAPIConfig", req.Msg, resp)
}

// UpdateAPIConfig 更新指定的API配置
func (u *UserServer) UpdateAPIConfig(ctx context.Context, req *connect.Request[pb.UpdateAPIConfigRequest]) (
	*connect.Response[pb.UpdateAPIConfigResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdateAPIConfig", req.Msg, (*pb.UpdateAPIConfigResponse)(nil), err)
	}

	resp, err := u.user.UpdateAPIConfig(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateAPIConfig", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateAPIConfig", req.Msg, resp)
}

// DeleteAPIConfig 删除指定的API配置
func (u *UserServer) DeleteAPIConfig(ctx context.Context, req *connect.Request[pb.DeleteAPIConfigRequest]) (
	*connect.Response[pb.DeleteAPIConfigResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "DeleteAPIConfig", req.Msg, (*pb.DeleteAPIConfigResponse)(nil), err)
	}

	resp, err := u.user.DeleteAPIConfig(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeleteAPIConfig", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeleteAPIConfig", req.Msg, resp)
}

// SetActiveAPIConfig 设置当前激活的API配置
func (u *UserServer) SetActiveAPIConfig(ctx context.Context, req *connect.Request[pb.SetActiveAPIConfigRequest]) (
	*connect.Response[pb.SetActiveAPIConfigResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "SetActiveAPIConfig", req.Msg, (*pb.SetActiveAPIConfigResponse)(nil), err)
	}

	resp, err := u.user.SetActiveAPIConfig(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "SetActiveAPIConfig", req.Msg, resp, err)
	}
	return doResponse(ctx, "SetActiveAPIConfig", req.Msg, resp)
}

// TestAPIConfig 测试API配置的连通性
func (u *UserServer) TestAPIConfig(ctx context.Context, req *connect.Request[pb.TestAPIConfigRequest]) (
	*connect.Response[pb.TestAPIConfigResponse], error) {
	resp, err := u.user.TestAPIConfig(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "TestAPIConfig", req.Msg, resp, err)
	}
	return doResponse(ctx, "TestAPIConfig", req.Msg, resp)
}
