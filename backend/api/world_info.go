package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/world_info"
)

type WorldInfoServer struct {
	worldInfo world_info.WorldInfo
}

func NewWorldInfoServer() *WorldInfoServer {
	return &WorldInfoServer{
		worldInfo: world_info.NewWorldInfo(),
	}
}

func (w *WorldInfoServer) ListWorldInfos(ctx context.Context, req *connect.Request[pb.ListWorldInfosRequest]) (
	*connect.Response[pb.ListWorldInfosResponse], error) {
	resp, err := w.worldInfo.ListWorldInfos(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoServer) GetWorldInfo(ctx context.Context, req *connect.Request[pb.GetWorldInfoRequest]) (
	*connect.Response[pb.GetWorldInfoResponse], error) {
	resp, err := w.worldInfo.GetWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoServer) CreateWorldInfo(ctx context.Context, req *connect.Request[pb.CreateWorldInfoRequest]) (
	*connect.Response[pb.CreateWorldInfoResponse], error) {
	resp, err := w.worldInfo.CreateWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoServer) UpdateWorldInfo(ctx context.Context, req *connect.Request[pb.UpdateWorldInfoRequest]) (
	*connect.Response[pb.UpdateWorldInfoResponse], error) {
	resp, err := w.worldInfo.UpdateWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoServer) DeleteWorldInfo(ctx context.Context, req *connect.Request[pb.DeleteWorldInfoRequest]) (
	*connect.Response[pb.DeleteWorldInfoResponse], error) {
	resp, err := w.worldInfo.DeleteWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoServer) ImportWorldInfo(ctx context.Context, req *connect.Request[pb.ImportWorldInfoRequest]) (
	*connect.Response[pb.ImportWorldInfoResponse], error) {
	resp, err := w.worldInfo.ImportWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoServer) ExportWorldInfo(ctx context.Context, req *connect.Request[pb.ExportWorldInfoRequest]) (
	*connect.Response[pb.ExportWorldInfoResponse], error) {
	resp, err := w.worldInfo.ExportWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
