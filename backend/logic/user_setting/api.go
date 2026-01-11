package user_setting

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type UserSetting interface {
	GetUserSetting(ctx context.Context, req *pb.GetUserSettingRequest) (*pb.GetUserSettingResponse, error)
	UpdateUserSetting(ctx context.Context, req *pb.UpdateUserSettingRequest) (*pb.UpdateUserSettingResponse, error)
}

func NewUserSetting() UserSetting {
	return newUserSetting()
}
