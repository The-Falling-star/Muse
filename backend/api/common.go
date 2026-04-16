package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/common"
)

// CommonServer 公共服务
type CommonServer struct {
	common common.Common
}

// NewCommonServer 创建一个新的CommonServer实例
func NewCommonServer() *CommonServer {
	return &CommonServer{
		common: common.NewCommon(),
	}
}

// GetPublicConfig 获取公共配置（无需认证）
func (c *CommonServer) GetPublicConfig(ctx context.Context, req *connect.Request[pb.GetPublicConfigRequest]) (
	*connect.Response[pb.GetPublicConfigResponse], error) {
	resp, err := c.common.GetPublicConfig(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetPublicConfig", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetPublicConfig", req.Msg, resp)
}
