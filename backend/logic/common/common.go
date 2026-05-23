package common

import (
	"context"

	"github.com/ling/muse/config"
	pb "github.com/ling/muse/gen/muse"
)

type commonImpl struct{}

func newCommon() *commonImpl {
	return &commonImpl{}
}

// GetPublicConfig 获取公共配置（无需认证）
func (c *commonImpl) GetPublicConfig(ctx context.Context, req *pb.GetPublicConfigRequest) (*pb.GetPublicConfigResponse, error) {
	cfg := config.Get()

	// 构建候选模型列表
	candidateModels := make([]*pb.CandidateModels, 0, len(cfg.CandidateModels))
	for providerName, models := range cfg.CandidateModels {
		candidateModels = append(candidateModels, &pb.CandidateModels{
			Provider: pb.APIProvider(pb.APIProvider_value[providerName]),
			Models:   models,
		})
	}

	return &pb.GetPublicConfigResponse{
		SkipAuth:        cfg.Auth.SkipAuth,
		CandidateModels: candidateModels,
	}, nil
}
