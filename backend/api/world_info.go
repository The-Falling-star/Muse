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
func (w *WorldInfoServer) ListWorldInfos(ctx context.Context, req *connect.Request[pb.ListWorldInfosReq]) (
	*connect.Response[pb.ListWorldInfosRsp], error) {
	resp, err := w.worldInfo.ListWorldInfos(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ListWorldInfos", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListWorldInfos", req.Msg, resp)
}

// GetWorldInfo 获取指定世界书详情
func (w *WorldInfoServer) GetWorldInfo(ctx context.Context, req *connect.Request[pb.GetWorldInfoReq]) (
	*connect.Response[pb.GetWorldInfoRsp], error) {
	resp, err := w.worldInfo.GetWorldInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetWorldInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetWorldInfo", req.Msg, resp)
}

// CreateWorldInfo 创建新的世界书
func (w *WorldInfoServer) CreateWorldInfo(ctx context.Context, req *connect.Request[pb.CreateWorldInfoReq]) (
	*connect.Response[pb.CreateWorldInfoRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "CreateWorldInfo", req.Msg, (*pb.CreateWorldInfoRsp)(nil), err)
	}

	resp, err := w.worldInfo.CreateWorldInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "CreateWorldInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "CreateWorldInfo", req.Msg, resp)
}

// UpdateWorldInfo 更新指定世界书
func (w *WorldInfoServer) UpdateWorldInfo(ctx context.Context, req *connect.Request[pb.UpdateWorldInfoReq]) (
	*connect.Response[pb.UpdateWorldInfoRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdateWorldInfo", req.Msg, (*pb.UpdateWorldInfoRsp)(nil), err)
	}

	resp, err := w.worldInfo.UpdateWorldInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateWorldInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateWorldInfo", req.Msg, resp)
}

// DeleteWorldInfo 删除指定世界书
func (w *WorldInfoServer) DeleteWorldInfo(ctx context.Context, req *connect.Request[pb.DeleteWorldInfoReq]) (
	*connect.Response[pb.DeleteWorldInfoRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "DeleteWorldInfo", req.Msg, (*pb.DeleteWorldInfoRsp)(nil), err)
	}

	resp, err := w.worldInfo.DeleteWorldInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeleteWorldInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeleteWorldInfo", req.Msg, resp)
}

// ImportWorldInfo 导入世界书
func (w *WorldInfoServer) ImportWorldInfo(ctx context.Context, req *connect.Request[pb.ImportWorldInfoReq]) (
	*connect.Response[pb.ImportWorldInfoRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "ImportWorldInfo", req.Msg, (*pb.ImportWorldInfoRsp)(nil), err)
	}

	resp, err := w.worldInfo.ImportWorldInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ImportWorldInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "ImportWorldInfo", req.Msg, resp)
}

// ExportWorldInfo 导出世界书
func (w *WorldInfoServer) ExportWorldInfo(ctx context.Context, req *connect.Request[pb.ExportWorldInfoReq]) (
	*connect.Response[pb.ExportWorldInfoRsp], error) {
	resp, err := w.worldInfo.ExportWorldInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ExportWorldInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "ExportWorldInfo", req.Msg, resp)
}

// ListWorldInfoEntries 获取世界书条目列表
func (w *WorldInfoServer) ListWorldInfoEntries(ctx context.Context,
	req *connect.Request[pb.ListWorldInfoEntriesReq]) (
	*connect.Response[pb.ListWorldInfoEntriesRsp], error) {
	resp, err := w.worldInfo.ListWorldInfoEntries(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ListWorldInfoEntries", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListWorldInfoEntries", req.Msg, resp)
}

// AddWorldInfoEntry 添加新的世界书条目
func (w *WorldInfoServer) AddWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.AddWorldInfoEntryReq]) (
	*connect.Response[pb.AddWorldInfoEntryRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "AddWorldInfoEntry", req.Msg, (*pb.AddWorldInfoEntryRsp)(nil), err)
	}

	resp, err := w.worldInfo.AddWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "AddWorldInfoEntry", req.Msg, resp, err)
	}
	return doResponse(ctx, "AddWorldInfoEntry", req.Msg, resp)
}

// UpdateWorldInfoEntry 更新指定世界书条目
func (w *WorldInfoServer) UpdateWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.UpdateWorldInfoEntryReq]) (
	*connect.Response[pb.UpdateWorldInfoEntryRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdateWorldInfoEntry", req.Msg, (*pb.UpdateWorldInfoEntryRsp)(nil), err)
	}

	resp, err := w.worldInfo.UpdateWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateWorldInfoEntry", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateWorldInfoEntry", req.Msg, resp)
}

// DeleteWorldInfoEntry 删除指定世界书条目
func (w *WorldInfoServer) DeleteWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.DeleteWorldInfoEntryReq]) (
	*connect.Response[pb.DeleteWorldInfoEntryRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "DeleteWorldInfoEntry", req.Msg, (*pb.DeleteWorldInfoEntryRsp)(nil), err)
	}

	resp, err := w.worldInfo.DeleteWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeleteWorldInfoEntry", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeleteWorldInfoEntry", req.Msg, resp)
}

// UpdateWorldInfoEntriesOrder 更新世界书条目的排序
func (w *WorldInfoServer) UpdateWorldInfoEntriesOrder(ctx context.Context,
	req *connect.Request[pb.UpdateWorldInfoEntriesOrderReq]) (
	*connect.Response[pb.UpdateWorldInfoEntriesOrderRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdateWorldInfoEntriesOrder", req.Msg, (*pb.UpdateWorldInfoEntriesOrderRsp)(nil), err)
	}

	resp, err := w.worldInfo.UpdateWorldInfoEntriesOrder(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateWorldInfoEntriesOrder", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateWorldInfoEntriesOrder", req.Msg, resp)
}
