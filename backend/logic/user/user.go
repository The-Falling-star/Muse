package user

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type userImpl struct {
}

func newUser() *userImpl {
	return &userImpl{}
}

func (u *userImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) GetCurrentUser(ctx context.Context, req *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) ListPersonas(ctx context.Context, req *pb.ListPersonasRequest) (*pb.ListPersonasResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) GetPersona(ctx context.Context, req *pb.GetPersonaRequest) (*pb.GetPersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) CreatePersona(ctx context.Context, req *pb.CreatePersonaRequest) (*pb.CreatePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) UpdatePersona(ctx context.Context, req *pb.UpdatePersonaRequest) (*pb.UpdatePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) DeletePersona(ctx context.Context, req *pb.DeletePersonaRequest) (*pb.DeletePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) SetActivePersona(ctx context.Context, req *pb.SetActivePersonaRequest) (*pb.SetActivePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) GetUserSetting(ctx context.Context, req *pb.GetUserSettingRequest) (*pb.GetUserSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) UpdateUserSetting(ctx context.Context, req *pb.UpdateUserSettingRequest) (*pb.UpdateUserSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) ListAPIConfigs(ctx context.Context, req *pb.ListAPIConfigsRequest) (*pb.ListAPIConfigsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) GetAPIConfig(ctx context.Context, req *pb.GetAPIConfigRequest) (*pb.GetAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) CreateAPIConfig(ctx context.Context, req *pb.CreateAPIConfigRequest) (*pb.CreateAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) UpdateAPIConfig(ctx context.Context, req *pb.UpdateAPIConfigRequest) (*pb.UpdateAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) DeleteAPIConfig(ctx context.Context, req *pb.DeleteAPIConfigRequest) (*pb.DeleteAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) SetActiveAPIConfig(ctx context.Context, req *pb.SetActiveAPIConfigRequest) (*pb.SetActiveAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userImpl) TestAPIConfig(ctx context.Context, req *pb.TestAPIConfigRequest) (*pb.TestAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}
