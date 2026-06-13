package chat

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// SendMessageStream 用于流式发送消息的回调接口
type SendMessageStream interface {
	Send(*pb.SendMessageRsp) error
}

// RegenerateMessageStream 用于流式重新生成消息的回调接口
type RegenerateMessageStream interface {
	Send(*pb.RegenerateMessageRsp) error
}

// Chat 定义了聊天服务的接口
type Chat interface {
	// ListChatSessions 获取聊天会话列表
	ListChatSessions(ctx context.Context, req *pb.ListChatSessionsReq) (*pb.ListChatSessionsRsp, error)
	// GetChatSession 获取指定聊天会话详情
	GetChatSession(ctx context.Context, req *pb.GetChatSessionReq) (*pb.GetChatSessionRsp, error)
	// CreateChatSession 创建新的聊天会话
	CreateChatSession(ctx context.Context, req *pb.CreateChatSessionReq) (*pb.CreateChatSessionRsp, error)
	// UpdateChatSession 更新指定聊天会话
	UpdateChatSession(ctx context.Context, req *pb.UpdateChatSessionReq) (*pb.UpdateChatSessionRsp, error)
	// DeleteChatSession 删除指定聊天会话
	DeleteChatSession(ctx context.Context, req *pb.DeleteChatSessionReq) (*pb.DeleteChatSessionRsp, error)
	// SendMessage 发送消息（流式响应）
	SendMessage(ctx context.Context, req *pb.SendMessageReq, stream SendMessageStream) error
	// RegenerateMessage 重新生成消息（流式响应）
	RegenerateMessage(ctx context.Context, req *pb.RegenerateMessageReq, stream RegenerateMessageStream) error
	// EditMessage 编辑消息
	EditMessage(ctx context.Context, req *pb.EditMessageReq) (*pb.EditMessageRsp, error)
	// DeleteMessage 删除消息
	DeleteMessage(ctx context.Context, req *pb.DeleteMessageReq) (*pb.DeleteMessageRsp, error)
	// SwitchSwipe 切换消息滑动选项
	SwitchSwipe(ctx context.Context, req *pb.SwitchSwipeReq) (*pb.SwitchSwipeRsp, error)
	// GetCharLatestSession 获取角色的最新
	GetCharLatestSession(ctx context.Context, req *pb.GetCharLatestSessionReq) (*pb.GetCharLatestSessionRsp, error)
	// UpdateSessionTime 更新会话时间
	UpdateSessionTime(ctx context.Context, req *pb.UpdateSessionTimeReq) (*pb.UpdateSessionTimeRsp, error)
	// DeleteSwipe 删除消息swipe
	DeleteSwipe(ctx context.Context, req *pb.DeleteSwipeReq) (*pb.DeleteSwipeRsp, error)
	// ImportSession 导入会话
	ImportSession(ctx context.Context, req *pb.ImportSessionReq) (*pb.ImportSessionRsp, error)
}

// NewChat 创建一个新的Chat实例
func NewChat() Chat {
	return newChat()
}
