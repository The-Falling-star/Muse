package api_config

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type apiConfigImpl struct {
}

func newAPIConfig() *apiConfigImpl {
	return &apiConfigImpl{}
}

func (a *apiConfigImpl) ListAPIConfigs(ctx context.Context, req *pb.ListAPIConfigsRequest) (*pb.ListAPIConfigsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *apiConfigImpl) GetAPIConfig(ctx context.Context, req *pb.GetAPIConfigRequest) (*pb.GetAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *apiConfigImpl) CreateAPIConfig(ctx context.Context, req *pb.CreateAPIConfigRequest) (*pb.CreateAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *apiConfigImpl) UpdateAPIConfig(ctx context.Context, req *pb.UpdateAPIConfigRequest) (*pb.UpdateAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *apiConfigImpl) DeleteAPIConfig(ctx context.Context, req *pb.DeleteAPIConfigRequest) (*pb.DeleteAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *apiConfigImpl) SetActiveAPIConfig(ctx context.Context, req *pb.SetActiveAPIConfigRequest) (*pb.SetActiveAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *apiConfigImpl) TestAPIConfig(ctx context.Context, req *pb.TestAPIConfigRequest) (*pb.TestAPIConfigResponse, error) {
	//TODO implement me
	panic("implement me")
}
