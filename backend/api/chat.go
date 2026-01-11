package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/chat"
)

// ChatServer 聊天服务
type ChatServer struct {
	chat chat.Chat
}

// NewChatServer 创建一个新的ChatServer实例
func NewChatServer() *ChatServer {
	return &ChatServer{
		chat: chat.NewChat(),
	}
}

// ListChatSessions 获取聊天会话列表
func (c *ChatServer) ListChatSessions(ctx context.Context, req *connect.Request[pb.ListChatSessionsRequest]) (
	*connect.Response[pb.ListChatSessionsResponse], error) {
	resp, err := c.chat.ListChatSessions(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetChatSession 获取指定聊天会话详情
func (c *ChatServer) GetChatSession(ctx context.Context, req *connect.Request[pb.GetChatSessionRequest]) (
	*connect.Response[pb.GetChatSessionResponse], error) {
	resp, err := c.chat.GetChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// CreateChatSession 创建新的聊天会话
func (c *ChatServer) CreateChatSession(ctx context.Context, req *connect.Request[pb.CreateChatSessionRequest]) (
	*connect.Response[pb.CreateChatSessionResponse], error) {
	resp, err := c.chat.CreateChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateChatSession 更新指定聊天会话
func (c *ChatServer) UpdateChatSession(ctx context.Context, req *connect.Request[pb.UpdateChatSessionRequest]) (
	*connect.Response[pb.UpdateChatSessionResponse], error) {
	resp, err := c.chat.UpdateChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteChatSession 删除指定聊天会话
func (c *ChatServer) DeleteChatSession(ctx context.Context, req *connect.Request[pb.DeleteChatSessionRequest]) (
	*connect.Response[pb.DeleteChatSessionResponse], error) {
	resp, err := c.chat.DeleteChatSession(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// SendMessage 发送消息（流式响应）
func (c *ChatServer) SendMessage(ctx context.Context, req *connect.Request[pb.SendMessageRequest],
	stream *connect.ServerStream[pb.SendMessageResponse]) error {
	return c.chat.SendMessage(ctx, req.Msg, stream)
}

// RegenerateMessage 重新生成消息（流式响应）
func (c *ChatServer) RegenerateMessage(ctx context.Context, req *connect.Request[pb.RegenerateMessageRequest],
	stream *connect.ServerStream[pb.RegenerateMessageResponse]) error {
	return c.chat.RegenerateMessage(ctx, req.Msg, stream)
}

// EditMessage 编辑消息
func (c *ChatServer) EditMessage(ctx context.Context, req *connect.Request[pb.EditMessageRequest]) (
	*connect.Response[pb.EditMessageResponse], error) {
	resp, err := c.chat.EditMessage(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteMessage 删除消息
func (c *ChatServer) DeleteMessage(ctx context.Context, req *connect.Request[pb.DeleteMessageRequest]) (
	*connect.Response[pb.DeleteMessageResponse], error) {
	resp, err := c.chat.DeleteMessage(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// SwitchSwipe 切换消息滑动选项
func (c *ChatServer) SwitchSwipe(ctx context.Context, req *connect.Request[pb.SwitchSwipeRequest]) (
	*connect.Response[pb.SwitchSwipeResponse], error) {
	resp, err := c.chat.SwitchSwipe(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
