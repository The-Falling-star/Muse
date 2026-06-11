package common

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// Common 公共服务接口
type Common interface {
	// GetPublicConfig 获取公共配置（无需认证）
	GetPublicConfig(ctx context.Context, req *pb.GetPublicConfigReq) (*pb.GetPublicConfigRsp, error)
}

// NewCommon 创建新的Common实例
func NewCommon() Common {
	return newCommon()
}
