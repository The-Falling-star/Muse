package api_config

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type APIConfig interface {
	ListAPIConfigs(ctx context.Context, req *pb.ListAPIConfigsRequest) (*pb.ListAPIConfigsResponse, error)
	GetAPIConfig(ctx context.Context, req *pb.GetAPIConfigRequest) (*pb.GetAPIConfigResponse, error)
	CreateAPIConfig(ctx context.Context, req *pb.CreateAPIConfigRequest) (*pb.CreateAPIConfigResponse, error)
	UpdateAPIConfig(ctx context.Context, req *pb.UpdateAPIConfigRequest) (*pb.UpdateAPIConfigResponse, error)
	DeleteAPIConfig(ctx context.Context, req *pb.DeleteAPIConfigRequest) (*pb.DeleteAPIConfigResponse, error)
	SetActiveAPIConfig(ctx context.Context, req *pb.SetActiveAPIConfigRequest) (*pb.SetActiveAPIConfigResponse, error)
	TestAPIConfig(ctx context.Context, req *pb.TestAPIConfigRequest) (*pb.TestAPIConfigResponse, error)
}

func NewAPIConfig() APIConfig {
	return newAPIConfig()
}
