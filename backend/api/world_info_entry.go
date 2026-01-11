package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/world_info_entry"
)

type WorldInfoEntryServer struct {
	worldInfoEntry world_info_entry.WorldInfoEntry
}

func NewWorldInfoEntryServer() *WorldInfoEntryServer {
	return &WorldInfoEntryServer{
		worldInfoEntry: world_info_entry.NewWorldInfoEntry(),
	}
}

func (w *WorldInfoEntryServer) ListWorldInfoEntries(ctx context.Context,
	req *connect.Request[pb.ListWorldInfoEntriesRequest]) (
	*connect.Response[pb.ListWorldInfoEntriesResponse], error) {
	resp, err := w.worldInfoEntry.ListWorldInfoEntries(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoEntryServer) AddWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.AddWorldInfoEntryRequest]) (
	*connect.Response[pb.AddWorldInfoEntryResponse], error) {
	resp, err := w.worldInfoEntry.AddWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoEntryServer) UpdateWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.UpdateWorldInfoEntryRequest]) (
	*connect.Response[pb.UpdateWorldInfoEntryResponse], error) {
	resp, err := w.worldInfoEntry.UpdateWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoEntryServer) DeleteWorldInfoEntry(ctx context.Context,
	req *connect.Request[pb.DeleteWorldInfoEntryRequest]) (
	*connect.Response[pb.DeleteWorldInfoEntryResponse], error) {
	resp, err := w.worldInfoEntry.DeleteWorldInfoEntry(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (w *WorldInfoEntryServer) UpdateWorldInfoEntriesOrder(ctx context.Context,
	req *connect.Request[pb.UpdateWorldInfoEntriesOrderRequest]) (
	*connect.Response[pb.UpdateWorldInfoEntriesOrderResponse], error) {
	resp, err := w.worldInfoEntry.UpdateWorldInfoEntriesOrder(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
