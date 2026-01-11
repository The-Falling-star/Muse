package chat

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// SendMessageStream 用于流式发送消息的回调接口
type SendMessageStream interface {
	Send(*pb.SendMessageResponse) error
}

// RegenerateMessageStream 用于流式重新生成消息的回调接口
type RegenerateMessageStream interface {
	Send(*pb.RegenerateMessageResponse) error
}

// Chat 定义了聊天服务的接口
type Chat interface {
	// ListChatSessions 获取聊天会话列表
	ListChatSessions(ctx context.Context, req *pb.ListChatSessionsRequest) (*pb.ListChatSessionsResponse, error)
	// GetChatSession 获取指定聊天会话详情
	GetChatSession(ctx context.Context, req *pb.GetChatSessionRequest) (*pb.GetChatSessionResponse, error)
	// CreateChatSession 创建新的聊天会话
	CreateChatSession(ctx context.Context, req *pb.CreateChatSessionRequest) (*pb.CreateChatSessionResponse, error)
	// UpdateChatSession 更新指定聊天会话
	UpdateChatSession(ctx context.Context, req *pb.UpdateChatSessionRequest) (*pb.UpdateChatSessionResponse, error)
	// DeleteChatSession 删除指定聊天会话
	DeleteChatSession(ctx context.Context, req *pb.DeleteChatSessionRequest) (*pb.DeleteChatSessionResponse, error)
	// SendMessage 发送消息（流式响应）
	SendMessage(ctx context.Context, req *pb.SendMessageRequest, stream SendMessageStream) error
	// RegenerateMessage 重新生成消息（流式响应）
	RegenerateMessage(ctx context.Context, req *pb.RegenerateMessageRequest, stream RegenerateMessageStream) error
	// EditMessage 编辑消息
	EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error)
	// DeleteMessage 删除消息
	DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error)
	// SwitchSwipe 切换消息滑动选项
	SwitchSwipe(ctx context.Context, req *pb.SwitchSwipeRequest) (*pb.SwitchSwipeResponse, error)
}

// NewChat 创建一个新的Chat实例
func NewChat() Chat {
	return newChat()
}
