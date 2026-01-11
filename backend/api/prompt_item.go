package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/prompt_item"
)

type PromptItemServer struct {
	promptItem prompt_item.PromptItem
}

func NewPromptItemServer() *PromptItemServer {
	return &PromptItemServer{
		promptItem: prompt_item.NewPromptItem(),
	}
}

func (p *PromptItemServer) ListPromptItems(ctx context.Context, req *connect.Request[pb.ListPromptItemsRequest]) (
	*connect.Response[pb.ListPromptItemsResponse], error) {
	resp, err := p.promptItem.ListPromptItems(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PromptItemServer) AddPromptItem(ctx context.Context, req *connect.Request[pb.AddPromptItemRequest]) (
	*connect.Response[pb.AddPromptItemResponse], error) {
	resp, err := p.promptItem.AddPromptItem(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PromptItemServer) UpdatePromptItem(ctx context.Context, req *connect.Request[pb.UpdatePromptItemRequest]) (
	*connect.Response[pb.UpdatePromptItemResponse], error) {
	resp, err := p.promptItem.UpdatePromptItem(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PromptItemServer) DeletePromptItem(ctx context.Context, req *connect.Request[pb.DeletePromptItemRequest]) (
	*connect.Response[pb.DeletePromptItemResponse], error) {
	resp, err := p.promptItem.DeletePromptItem(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PromptItemServer) UpdatePromptItemsOrder(ctx context.Context,
	req *connect.Request[pb.UpdatePromptItemsOrderRequest]) (
	*connect.Response[pb.UpdatePromptItemsOrderResponse], error) {
	resp, err := p.promptItem.UpdatePromptItemsOrder(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
