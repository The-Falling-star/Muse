package world_info

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type WorldInfo interface {
	ListWorldInfos(ctx context.Context, req *pb.ListWorldInfosRequest) (*pb.ListWorldInfosResponse, error)
	GetWorldInfo(ctx context.Context, req *pb.GetWorldInfoRequest) (*pb.GetWorldInfoResponse, error)
	CreateWorldInfo(ctx context.Context, req *pb.CreateWorldInfoRequest) (*pb.CreateWorldInfoResponse, error)
	UpdateWorldInfo(ctx context.Context, req *pb.UpdateWorldInfoRequest) (*pb.UpdateWorldInfoResponse, error)
	DeleteWorldInfo(ctx context.Context, req *pb.DeleteWorldInfoRequest) (*pb.DeleteWorldInfoResponse, error)
	ImportWorldInfo(ctx context.Context, req *pb.ImportWorldInfoRequest) (*pb.ImportWorldInfoResponse, error)
	ExportWorldInfo(ctx context.Context, req *pb.ExportWorldInfoRequest) (*pb.ExportWorldInfoResponse, error)
}

func NewWorldInfo() WorldInfo {
	return newWorldInfo()
}
