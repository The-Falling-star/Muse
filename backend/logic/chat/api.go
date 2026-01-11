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

type Chat interface {
	ListChatSessions(ctx context.Context, req *pb.ListChatSessionsRequest) (*pb.ListChatSessionsResponse, error)
	GetChatSession(ctx context.Context, req *pb.GetChatSessionRequest) (*pb.GetChatSessionResponse, error)
	CreateChatSession(ctx context.Context, req *pb.CreateChatSessionRequest) (*pb.CreateChatSessionResponse, error)
	UpdateChatSession(ctx context.Context, req *pb.UpdateChatSessionRequest) (*pb.UpdateChatSessionResponse, error)
	DeleteChatSession(ctx context.Context, req *pb.DeleteChatSessionRequest) (*pb.DeleteChatSessionResponse, error)
	SendMessage(ctx context.Context, req *pb.SendMessageRequest, stream SendMessageStream) error
	RegenerateMessage(ctx context.Context, req *pb.RegenerateMessageRequest, stream RegenerateMessageStream) error
	EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error)
	DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error)
	SwitchSwipe(ctx context.Context, req *pb.SwitchSwipeRequest) (*pb.SwitchSwipeResponse, error)
}

func NewChat() Chat {
	return newChat()
}
