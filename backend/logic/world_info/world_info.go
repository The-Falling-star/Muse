package world_info

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type worldInfoImpl struct {
}

func newWorldInfo() *worldInfoImpl {
	return &worldInfoImpl{}
}

func (w *worldInfoImpl) ListWorldInfos(ctx context.Context, req *pb.ListWorldInfosRequest) (*pb.ListWorldInfosResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoImpl) GetWorldInfo(ctx context.Context, req *pb.GetWorldInfoRequest) (*pb.GetWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoImpl) CreateWorldInfo(ctx context.Context, req *pb.CreateWorldInfoRequest) (*pb.CreateWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoImpl) UpdateWorldInfo(ctx context.Context, req *pb.UpdateWorldInfoRequest) (*pb.UpdateWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoImpl) DeleteWorldInfo(ctx context.Context, req *pb.DeleteWorldInfoRequest) (*pb.DeleteWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoImpl) ImportWorldInfo(ctx context.Context, req *pb.ImportWorldInfoRequest) (*pb.ImportWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worldInfoImpl) ExportWorldInfo(ctx context.Context, req *pb.ExportWorldInfoRequest) (*pb.ExportWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}
