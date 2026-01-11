package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/character"
)

// CharacterServer 角色卡服务
type CharacterServer struct {
	character character.Character
}

// NewCharacterServer 创建一个新的CharacterServer实例
func NewCharacterServer() *CharacterServer {
	return &CharacterServer{
		character: character.NewCharacter(),
	}
}

// ListCharacters 获取角色列表
func (c *CharacterServer) ListCharacters(ctx context.Context, req *connect.Request[pb.ListCharactersRequest]) (
	*connect.Response[pb.ListCharactersResponse], error) {
	resp, err := c.character.ListCharacters(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetCharacter 获取指定角色详情
func (c *CharacterServer) GetCharacter(ctx context.Context, req *connect.Request[pb.GetCharacterRequest]) (
	*connect.Response[pb.GetCharacterResponse], error) {
	resp, err := c.character.GetCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// CreateCharacter 创建新角色
func (c *CharacterServer) CreateCharacter(ctx context.Context, req *connect.Request[pb.CreateCharacterRequest]) (
	*connect.Response[pb.CreateCharacterResponse], error) {
	resp, err := c.character.CreateCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateCharacter 更新指定角色
func (c *CharacterServer) UpdateCharacter(ctx context.Context, req *connect.Request[pb.UpdateCharacterRequest]) (
	*connect.Response[pb.UpdateCharacterResponse], error) {
	resp, err := c.character.UpdateCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteCharacter 删除指定角色
func (c *CharacterServer) DeleteCharacter(ctx context.Context, req *connect.Request[pb.DeleteCharacterRequest]) (
	*connect.Response[pb.DeleteCharacterResponse], error) {
	resp, err := c.character.DeleteCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ImportCharacter 导入角色
func (c *CharacterServer) ImportCharacter(ctx context.Context, req *connect.Request[pb.ImportCharacterRequest]) (
	*connect.Response[pb.ImportCharacterResponse], error) {
	resp, err := c.character.ImportCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ExportCharacter 导出角色
func (c *CharacterServer) ExportCharacter(ctx context.Context, req *connect.Request[pb.ExportCharacterRequest]) (
	*connect.Response[pb.ExportCharacterResponse], error) {
	resp, err := c.character.ExportCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// RestoreCharacterWorldInfo 恢复角色的世界信息
func (c *CharacterServer) RestoreCharacterWorldInfo(ctx context.Context, req *connect.Request[pb.RestoreCharacterWorldInfoRequest]) (
	*connect.Response[pb.RestoreCharacterWorldInfoResponse], error) {
	resp, err := c.character.RestoreCharacterWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
