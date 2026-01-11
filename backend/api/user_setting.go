package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/user_setting"
)

type UserSettingServer struct {
	userSetting user_setting.UserSetting
}

func NewUserSettingServer() *UserSettingServer {
	return &UserSettingServer{
		userSetting: user_setting.NewUserSetting(),
	}
}

func (u *UserSettingServer) GetUserSetting(ctx context.Context, req *connect.Request[pb.GetUserSettingRequest]) (
	*connect.Response[pb.GetUserSettingResponse], error) {
	resp, err := u.userSetting.GetUserSetting(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (u *UserSettingServer) UpdateUserSetting(ctx context.Context, req *connect.Request[pb.UpdateUserSettingRequest]) (
	*connect.Response[pb.UpdateUserSettingResponse], error) {
	resp, err := u.userSetting.UpdateUserSetting(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
