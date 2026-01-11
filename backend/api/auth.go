package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/auth"
)

type AuthServer struct {
	auth auth.Auth
}

func NewAuthServer() *AuthServer {
	return &AuthServer{
		auth: auth.NewAuth(),
	}
}

func (a *AuthServer) Register(ctx context.Context, req *connect.Request[pb.RegisterRequest]) (
	*connect.Response[pb.RegisterResponse], error) {
	resp, err := a.auth.Register(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *AuthServer) Login(ctx context.Context, req *connect.Request[pb.LoginRequest]) (
	*connect.Response[pb.LoginResponse], error) {
	resp, err := a.auth.Login(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *AuthServer) GetCurrentUser(ctx context.Context, req *connect.Request[pb.GetCurrentUserRequest]) (
	*connect.Response[pb.GetCurrentUserResponse], error) {
	resp, err := a.auth.GetCurrentUser(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *AuthServer) ChangePassword(ctx context.Context, req *connect.Request[pb.ChangePasswordRequest]) (
	*connect.Response[pb.ChangePasswordResponse], error) {
	resp, err := a.auth.ChangePassword(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
