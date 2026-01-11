package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/chat"
)

type ChatServer struct {
	chat chat.Chat
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		chat: chat.NewChat(),
	}
}

func (c *ChatServer) ListChatSessions(ctx context.Context, req *connect.Request[pb.ListChatSessionsRequest]) (
	*connect.Response[pb.ListChatSessionsResponse], error) {
	resp, err := c.chat.ListChatSessions(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *ChatServer) GetChatSession(ctx context.Context, req *connect.Request[pb.GetChatSessionRequest]) (
	*connect.Response[pb.GetChatSessionResponse], error) {
	resp, err := c.chat.GetChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *ChatServer) CreateChatSession(ctx context.Context, req *connect.Request[pb.CreateChatSessionRequest]) (
	*connect.Response[pb.CreateChatSessionResponse], error) {
	resp, err := c.chat.CreateChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *ChatServer) UpdateChatSession(ctx context.Context, req *connect.Request[pb.UpdateChatSessionRequest]) (
	*connect.Response[pb.UpdateChatSessionResponse], error) {
	resp, err := c.chat.UpdateChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *ChatServer) DeleteChatSession(ctx context.Context, req *connect.Request[pb.DeleteChatSessionRequest]) (
	*connect.Response[pb.DeleteChatSessionResponse], error) {
	resp, err := c.chat.DeleteChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *ChatServer) SendMessage(ctx context.Context, req *connect.Request[pb.SendMessageRequest],
	stream *connect.ServerStream[pb.SendMessageResponse]) error {
	return c.chat.SendMessage(ctx, req.Msg, stream)
}

func (c *ChatServer) RegenerateMessage(ctx context.Context, req *connect.Request[pb.RegenerateMessageRequest],
	stream *connect.ServerStream[pb.RegenerateMessageResponse]) error {
	return c.chat.RegenerateMessage(ctx, req.Msg, stream)
}

func (c *ChatServer) EditMessage(ctx context.Context, req *connect.Request[pb.EditMessageRequest]) (
	*connect.Response[pb.EditMessageResponse], error) {
	resp, err := c.chat.EditMessage(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *ChatServer) DeleteMessage(ctx context.Context, req *connect.Request[pb.DeleteMessageRequest]) (
	*connect.Response[pb.DeleteMessageResponse], error) {
	resp, err := c.chat.DeleteMessage(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *ChatServer) SwitchSwipe(ctx context.Context, req *connect.Request[pb.SwitchSwipeRequest]) (
	*connect.Response[pb.SwitchSwipeResponse], error) {
	resp, err := c.chat.SwitchSwipe(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
