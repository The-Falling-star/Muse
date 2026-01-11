package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/world_info"
)

// WorldInfoServer 世界书服务（包含世界书和世界书条目的管理）
type WorldInfoServer struct {
	worldInfo world_info.WorldInfo
}

// NewWorldInfoServer 创建一个新的WorldInfoServer实例
func NewWorldInfoServer() *WorldInfoServer {
	return &WorldInfoServer{
		worldInfo: world_info.NewWorldInfo(),
	}
}

// ListWorldInfos 获取世界书列表
func (w *WorldInfoServer) ListWorldInfos(ctx context.Context, req *connect.Request[pb.ListWorldInfosRequest]) (
	*connect.Response[pb.ListWorldInfosResponse], error) {
	resp, err := w.worldInfo.ListWorldInfos(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetWorldInfo 获取指定世界书详情
func (w *WorldInfoServer) GetWorldInfo(ctx context.Context, req *connect.Request[pb.GetWorldInfoRequest]) (
	*connect.Response[pb.GetWorldInfoResponse], error) {
	resp, err := w.worldInfo.GetWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// CreateWorldInfo 创建新的世界书
func (w *WorldInfoServer) CreateWorldInfo(ctx context.Context, req *connect.Request[pb.CreateWorldInfoRequest]) (
	*connect.Response[pb.CreateWorldInfoResponse], error) {
	resp, err := w.worldInfo.CreateWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateWorldInfo 更新指定世界书
func (w *WorldInfoServer) UpdateWorldInfo(ctx context.Context, req *connect.Request[pb.UpdateWorldInfoRequest]) (
	*connect.Response[pb.UpdateWorldInfoResponse], error) {
	resp, err := w.worldInfo.UpdateWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteWorldInfo 删除指定世界书
func (w *WorldInfoServer) DeleteWorldInfo(ctx context.Context, req *connect.Request[pb.DeleteWorldInfoRequest]) (
	*connect.Response[pb.DeleteWorldInfoResponse], error) {
	resp, err := w.worldInfo.DeleteWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ImportWorldInfo 导入世界书
func (w *WorldInfoServer) ImportWorldInfo(ctx context.Context, req *connect.Request[pb.ImportWorldInfoRequest]) (
	*connect.Response[pb.ImportWorldInfoResponse], error) {
	resp, err := w.worldInfo.ImportWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ExportWorldInfo 导出世界书
func (w *WorldInfoServer) ExportWorldInfo(ctx context.Context, req *connect.Request[pb.ExportWorldInfoRequest]) (
	*connect.Response[pb.ExportWorldInfoResponse], error) {
	resp, err := w.worldInfo.ExportWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListWorldInfoEntries 获取世界书条目列表
func (w *WorldInfoServer) ListWorldInfoEntries(ctx context.Context,
	req *connect.Request[pb.ListWorldInfoEntriesRequest]) (
	*connect.Response[pb.ListWorldInfoEntriesResponse], error) {
	resp, err := w.worldInfo.ListWorldInfoEntries(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// AddWorldInfoEntry 添加新的世界书条目
func (w *WorldInfoServer) AddWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.AddWorldInfoEntryRequest]) (
	*connect.Response[pb.AddWorldInfoEntryResponse], error) {
	resp, err := w.worldInfo.AddWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateWorldInfoEntry 更新指定世界书条目
func (w *WorldInfoServer) UpdateWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.UpdateWorldInfoEntryRequest]) (
	*connect.Response[pb.UpdateWorldInfoEntryResponse], error) {
	resp, err := w.worldInfo.UpdateWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteWorldInfoEntry 删除指定世界书条目
func (w *WorldInfoServer) DeleteWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.DeleteWorldInfoEntryRequest]) (
	*connect.Response[pb.DeleteWorldInfoEntryResponse], error) {
	resp, err := w.worldInfo.DeleteWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateWorldInfoEntriesOrder 更新世界书条目的排序
func (w *WorldInfoServer) UpdateWorldInfoEntriesOrder(ctx context.Context,
	req *connect.Request[pb.UpdateWorldInfoEntriesOrderRequest]) (
	*connect.Response[pb.UpdateWorldInfoEntriesOrderResponse], error) {
	resp, err := w.worldInfo.UpdateWorldInfoEntriesOrder(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
