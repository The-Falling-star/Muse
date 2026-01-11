package user_setting

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type userSettingImpl struct {
}

func newUserSetting() *userSettingImpl {
	return &userSettingImpl{}
}

func (u *userSettingImpl) GetUserSetting(ctx context.Context, req *pb.GetUserSettingRequest) (*pb.GetUserSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userSettingImpl) UpdateUserSetting(ctx context.Context, req *pb.UpdateUserSettingRequest) (*pb.UpdateUserSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}
