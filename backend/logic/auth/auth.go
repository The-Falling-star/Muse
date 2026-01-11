package auth

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type authImpl struct {
}

func newAuth() *authImpl {
	return &authImpl{}
}

func (a *authImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *authImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *authImpl) GetCurrentUser(ctx context.Context, req *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *authImpl) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}
