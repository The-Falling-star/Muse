package world_info

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// WorldInfo 定义了世界信息服务的接口（包含世界书和世界书条目的管理）
type WorldInfo interface {
	// ListWorldInfos 获取世界信息列表
	ListWorldInfos(ctx context.Context, req *pb.ListWorldInfosRequest) (*pb.ListWorldInfosResponse, error)
	// GetWorldInfo 获取指定世界信息详情
	GetWorldInfo(ctx context.Context, req *pb.GetWorldInfoRequest) (*pb.GetWorldInfoResponse, error)
	// CreateWorldInfo 创建新的世界信息
	CreateWorldInfo(ctx context.Context, req *pb.CreateWorldInfoRequest) (*pb.CreateWorldInfoResponse, error)
	// UpdateWorldInfo 更新指定世界信息
	UpdateWorldInfo(ctx context.Context, req *pb.UpdateWorldInfoRequest) (*pb.UpdateWorldInfoResponse, error)
	// DeleteWorldInfo 删除指定世界信息
	DeleteWorldInfo(ctx context.Context, req *pb.DeleteWorldInfoRequest) (*pb.DeleteWorldInfoResponse, error)
	// ImportWorldInfo 导入世界信息
	ImportWorldInfo(ctx context.Context, req *pb.ImportWorldInfoRequest) (*pb.ImportWorldInfoResponse, error)
	// ExportWorldInfo 导出世界信息
	ExportWorldInfo(ctx context.Context, req *pb.ExportWorldInfoRequest) (*pb.ExportWorldInfoResponse, error)

	// ListWorldInfoEntries 获取世界信息条目列表
	ListWorldInfoEntries(ctx context.Context, req *pb.ListWorldInfoEntriesRequest) (*pb.ListWorldInfoEntriesResponse, error)
	// AddWorldInfoEntry 添加新的世界信息条目
	AddWorldInfoEntry(ctx context.Context, req *pb.AddWorldInfoEntryRequest) (*pb.AddWorldInfoEntryResponse, error)
	// UpdateWorldInfoEntry 更新指定世界信息条目
	UpdateWorldInfoEntry(ctx context.Context, req *pb.UpdateWorldInfoEntryRequest) (*pb.UpdateWorldInfoEntryResponse, error)
	// DeleteWorldInfoEntry 删除指定世界信息条目
	DeleteWorldInfoEntry(ctx context.Context, req *pb.DeleteWorldInfoEntryRequest) (*pb.DeleteWorldInfoEntryResponse, error)
	// UpdateWorldInfoEntriesOrder 更新世界信息条目的排序
	UpdateWorldInfoEntriesOrder(ctx context.Context, req *pb.UpdateWorldInfoEntriesOrderRequest) (*pb.UpdateWorldInfoEntriesOrderResponse, error)
}

// NewWorldInfo 创建一个新的WorldInfo实例
func NewWorldInfo() WorldInfo {
	return newWorldInfo()
}
