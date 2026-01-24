package chat

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/constrant"
	"github.com/ling/muse/common/convert"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/common/jwt"
	"github.com/ling/muse/entity"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/repo/database"
)

// 默认用户ID，待认证功能完成后替换
const defaultUserID = 1

type chatImpl struct {
	chatRepo *database.ChatRepo
}

func newChat() *chatImpl {
	return &chatImpl{
		chatRepo: &database.ChatRepo{},
	}
}

func (c *chatImpl) ListChatSessions(ctx context.Context, req *pb.ListChatSessionsRequest) (*pb.ListChatSessionsResponse, error) {
	// 获取分页参数
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	characterID := int(req.GetCharacterId())

	// 设置默认值
	if page <= 0 {
		page = constrant.DefaultPageNum
	}
	if pageSize <= 0 {
		pageSize = constrant.DefaultPageSize
	}
	if pageSize > constrant.MaxPageSize {
		pageSize = constrant.MaxPageSize
	}

	// 从数据库获取会话列表
	userId := jwt.GetUserId(ctx)
	sessions, total, err := c.chatRepo.ListSessions(userId, characterID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为pb格式
	pbSessions := make([]*pb.ChatSession, 0, len(sessions))
	for _, session := range sessions {
		pbSessions = append(pbSessions, convert.SessionEntityToPb(session))
	}

	return &pb.ListChatSessionsResponse{
		Sessions: pbSessions,
		Total:    int32(total),
	}, nil
}

func (c *chatImpl) GetChatSession(ctx context.Context, req *pb.GetChatSessionRequest) (
	*pb.GetChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidSessionID)
	}

	// 从数据库获取会话
	session, err := c.chatRepo.GetSessionByID(id, defaultUserID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errs.NewStandard(connect.CodeNotFound, errs.SessionNotFound)
	}

	return &pb.GetChatSessionResponse{
		Session: convert.SessionEntityToPb(session),
	}, nil
}

func (c *chatImpl) CreateChatSession(ctx context.Context, req *pb.CreateChatSessionRequest) (*pb.CreateChatSessionResponse, error) {
	// 参数校验
	characterID := int(req.GetCharacterId())
	if characterID <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidCharacterID)
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		name = "新会话"
	}

	// 构建会话实体
	session := &entity.ChatSession{
		UserID:      defaultUserID,
		CharacterID: characterID,
		Name:        name,
		Version:     1,
	}

	// 保存到数据库
	if err := c.chatRepo.CreateSession(session); err != nil {
		return nil, err
	}

	// 重新获取完整数据（包含关联）
	fullSession, err := c.chatRepo.GetSessionByID(session.ID, defaultUserID)
	if err != nil {
		return nil, err
	}

	return &pb.CreateChatSessionResponse{
		Session: convert.SessionEntityToPb(fullSession),
	}, nil
}

func (c *chatImpl) UpdateChatSession(ctx context.Context, req *pb.UpdateChatSessionRequest) (*pb.UpdateChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidSessionID)
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.EmptySessionName)
	}
	userId := jwt.GetUserId(ctx)
	// 更新数据库
	_, err := c.chatRepo.UpdateSession(int(req.GetId()), userId, int(req.GetVersion()), req.GetName())
	if err != nil {
		return nil, err
	}
	// 重新获取更新后的会话（包含关联数据）
	updatedSession, err := c.chatRepo.GetSessionByID(id, userId)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateChatSessionResponse{
		Session: convert.SessionEntityToPb(updatedSession),
	}, nil
}

func (c *chatImpl) DeleteChatSession(ctx context.Context, req *pb.DeleteChatSessionRequest) (*pb.DeleteChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, errs.NewStandard(connect.CodeInvalidArgument, errs.InvalidSessionID)
	}

	// 删除会话（会级联删除关联的消息和swipes）
	if err := c.chatRepo.DeleteSession(id, defaultUserID); err != nil {
		return nil, err
	}

	return &pb.DeleteChatSessionResponse{}, nil
}

func (c *chatImpl) SendMessage(ctx context.Context, req *pb.SendMessageRequest, stream SendMessageStream) error {
	//TODO implement me
	panic("implement me")
}

func (c *chatImpl) RegenerateMessage(ctx context.Context, req *pb.RegenerateMessageRequest, stream RegenerateMessageStream) error {
	//TODO implement me
	panic("implement me")
}

func (c *chatImpl) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *chatImpl) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *chatImpl) SwitchSwipe(ctx context.Context, req *pb.SwitchSwipeRequest) (*pb.SwitchSwipeResponse, error) {
	//TODO implement me
	panic("implement me")
}
