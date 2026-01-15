package chat

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/ling/muse/common/convert"
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
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 从数据库获取会话列表
	sessions, total, err := c.chatRepo.ListSessions(defaultUserID, characterID, page, pageSize)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

func (c *chatImpl) GetChatSession(ctx context.Context, req *pb.GetChatSessionRequest) (*pb.GetChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的会话ID"))
	}

	// 从数据库获取会话
	session, err := c.chatRepo.GetSessionByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if session == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("会话不存在"))
	}

	return &pb.GetChatSessionResponse{
		Session: convert.SessionEntityToPb(session),
	}, nil
}

func (c *chatImpl) CreateChatSession(ctx context.Context, req *pb.CreateChatSessionRequest) (*pb.CreateChatSessionResponse, error) {
	// 参数校验
	characterID := int(req.GetCharacterId())
	if characterID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("角色ID不能为空"))
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取完整数据（包含关联）
	fullSession, err := c.chatRepo.GetSessionByID(session.ID, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.CreateChatSessionResponse{
		Session: convert.SessionEntityToPb(fullSession),
	}, nil
}

func (c *chatImpl) UpdateChatSession(ctx context.Context, req *pb.UpdateChatSessionRequest) (*pb.UpdateChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的会话ID"))
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("会话名称不能为空"))
	}

	// 获取当前会话
	session, err := c.chatRepo.GetSessionByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if session == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("会话不存在"))
	}

	// 检查版本号
	if int64(session.Version) != req.GetVersion() {
		return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("数据已被修改，请刷新后重试"))
	}

	// 更新会话字段
	session.Name = name

	// 更新数据库
	if err := c.chatRepo.UpdateSession(session); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的会话（包含关联数据）
	updatedSession, err := c.chatRepo.GetSessionByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdateChatSessionResponse{
		Session: convert.SessionEntityToPb(updatedSession),
	}, nil
}

func (c *chatImpl) DeleteChatSession(ctx context.Context, req *pb.DeleteChatSessionRequest) (*pb.DeleteChatSessionResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的会话ID"))
	}

	// 删除会话（会级联删除关联的消息和swipes）
	if err := c.chatRepo.DeleteSession(id, defaultUserID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
