package api

import (
	"context"
	"time"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/chat"
	log "github.com/sirupsen/logrus"
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
func (c *ChatServer) ListChatSessions(ctx context.Context, req *connect.Request[pb.ListChatSessionsReq]) (
	*connect.Response[pb.ListChatSessionsRsp], error) {
	resp, err := c.chat.ListChatSessions(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ListChatSessions", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListChatSessions", req.Msg, resp)
}

// GetChatSession 获取指定聊天会话详情
func (c *ChatServer) GetChatSession(ctx context.Context, req *connect.Request[pb.GetChatSessionReq]) (
	*connect.Response[pb.GetChatSessionRsp], error) {
	resp, err := c.chat.GetChatSession(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetChatSession", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetChatSession", req.Msg, resp)
}

// CreateChatSession 创建新的聊天会话
func (c *ChatServer) CreateChatSession(ctx context.Context, req *connect.Request[pb.CreateChatSessionReq]) (
	*connect.Response[pb.CreateChatSessionRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := c.chat.CreateChatSession(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "CreateChatSession", req.Msg, resp, err)
	}
	return doResponse(ctx, "CreateChatSession", req.Msg, resp)
}

// UpdateChatSession 更新指定聊天会话
func (c *ChatServer) UpdateChatSession(ctx context.Context, req *connect.Request[pb.UpdateChatSessionReq]) (
	*connect.Response[pb.UpdateChatSessionRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := c.chat.UpdateChatSession(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateChatSession", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateChatSession", req.Msg, resp)
}

// DeleteChatSession 删除指定聊天会话
func (c *ChatServer) DeleteChatSession(ctx context.Context, req *connect.Request[pb.DeleteChatSessionReq]) (
	*connect.Response[pb.DeleteChatSessionRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := c.chat.DeleteChatSession(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeleteChatSession", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeleteChatSession", req.Msg, resp)
}

// SendMessage 发送消息（流式响应）
func (c *ChatServer) SendMessage(ctx context.Context, req *connect.Request[pb.SendMessageReq],
	stream *connect.ServerStream[pb.SendMessageRsp]) error {
	startTime := time.Now().UnixMilli()
	// 开启事务
	var (
		err error              = nil
		rsp *pb.SendMessageRsp = nil
	)

	ctx, err = beginTransaction(ctx)
	if err != nil {
		_, err = doResponseExp(ctx, "SendMessage", req.Msg, rsp, err)
		return err
	}
	if err = c.chat.SendMessage(ctx, req.Msg, stream); err != nil {
		_, err = doResponseExp(ctx, "SendMessage", req.Msg, rsp, err)
		return err
	}
	_, _ = doResponse(ctx, "SendMessage", req.Msg, rsp)
	log.Infof("调用耗时: %d ms", time.Now().UnixMilli()-startTime)
	return nil
}

// RegenerateMessage 重新生成消息（流式响应）
func (c *ChatServer) RegenerateMessage(ctx context.Context, req *connect.Request[pb.RegenerateMessageReq],
	stream *connect.ServerStream[pb.RegenerateMessageRsp]) error {
	// 开启事务
	var (
		err error                    = nil
		rsp *pb.RegenerateMessageRsp = nil
	)

	ctx, err = beginTransaction(ctx)
	if err != nil {
		_, err = doResponseExp(ctx, "RegenerateMessage", req.Msg, rsp, err)
		return err
	}
	if err = c.chat.RegenerateMessage(ctx, req.Msg, stream); err != nil {
		_, err = doResponseExp(ctx, "RegenerateMessage", req.Msg, rsp, err)
		return err
	}
	_, _ = doResponse(ctx, "RegenerateMessage", req.Msg, rsp)
	return nil
}

// EditMessage 编辑消息
func (c *ChatServer) EditMessage(ctx context.Context, req *connect.Request[pb.EditMessageReq]) (
	*connect.Response[pb.EditMessageRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := c.chat.EditMessage(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "EditMessage", req.Msg, resp, err)
	}
	return doResponse(ctx, "EditMessage", req.Msg, resp)
}

// DeleteMessage 删除消息
func (c *ChatServer) DeleteMessage(ctx context.Context, req *connect.Request[pb.DeleteMessageReq]) (
	*connect.Response[pb.DeleteMessageRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := c.chat.DeleteMessage(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeleteMessage", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeleteMessage", req.Msg, resp)
}

// SwitchSwipe 切换消息滑动选项
func (c *ChatServer) SwitchSwipe(ctx context.Context, req *connect.Request[pb.SwitchSwipeReq]) (
	*connect.Response[pb.SwitchSwipeRsp], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := c.chat.SwitchSwipe(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "SwitchSwipe", req.Msg, resp, err)
	}
	return doResponse(ctx, "SwitchSwipe", req.Msg, resp)
}

func (c *ChatServer) GetCharLatestSession(ctx context.Context, req *connect.Request[pb.GetCharLatestSessionReq]) (
	*connect.Response[pb.GetCharLatestSessionRsp], error) {
	resp, err := c.chat.GetCharLatestSession(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetCharLatestSession", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetCharLatestSession", req.Msg, resp)
}

func (c *ChatServer) UpdateSessionTime(ctx context.Context, req *connect.Request[pb.UpdateSessionTimeReq]) (
	*connect.Response[pb.UpdateSessionTimeRsp], error) {
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp, err := c.chat.UpdateSessionTime(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateSessionTime", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateSessionTime", req.Msg, resp)
}

func (c *ChatServer) DeleteSwipe(ctx context.Context, req *connect.Request[pb.DeleteSwipeReq]) (
	*connect.Response[pb.DeleteSwipeRsp], error) {
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp, err := c.chat.DeleteSwipe(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeleteSwipe", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeleteSwipe", req.Msg, resp)
}
