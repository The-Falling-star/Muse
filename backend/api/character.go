package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/character"
	log "github.com/sirupsen/logrus"
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
		return doResponseExp(ctx, "ListCharacters", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListCharacters", req.Msg, resp)
}

// GetCharacter 获取指定角色详情
func (c *CharacterServer) GetCharacter(ctx context.Context, req *connect.Request[pb.GetCharacterRequest]) (
	*connect.Response[pb.GetCharacterResponse], error) {
	resp, err := c.character.GetCharacter(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetCharacter", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetCharacter", req.Msg, resp)
}

// CreateCharacter 创建新角色
func (c *CharacterServer) CreateCharacter(ctx context.Context, req *connect.Request[pb.CreateCharacterRequest]) (
	*connect.Response[pb.CreateCharacterResponse], error) {
	resp, err := c.character.CreateCharacter(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "CreateCharacter", req.Msg, resp, err)
	}
	return doResponse(ctx, "CreateCharacter", req.Msg, resp)
}

// UpdateCharacter 更新指定角色
func (c *CharacterServer) UpdateCharacter(ctx context.Context, req *connect.Request[pb.UpdateCharacterRequest]) (
	*connect.Response[pb.UpdateCharacterResponse], error) {
	resp, err := c.character.UpdateCharacter(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdateCharacter", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdateCharacter", req.Msg, resp)
}

// DeleteCharacter 删除指定角色
func (c *CharacterServer) DeleteCharacter(ctx context.Context, req *connect.Request[pb.DeleteCharacterRequest]) (
	*connect.Response[pb.DeleteCharacterResponse], error) {
	resp, err := c.character.DeleteCharacter(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeleteCharacter", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeleteCharacter", req.Msg, resp)
}

// ImportCharacter 导入角色
func (c *CharacterServer) ImportCharacter(ctx context.Context, req *connect.Request[pb.ImportCharacterRequest]) (
	*connect.Response[pb.ImportCharacterResponse], error) {
	resp, err := c.character.ImportCharacter(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ImportCharacter", req.Msg, resp, err)
	}
	return doResponse(ctx, "ImportCharacter", req.Msg, resp)
}

// ExportCharacter 导出角色
func (c *CharacterServer) ExportCharacter(ctx context.Context, req *connect.Request[pb.ExportCharacterRequest]) (
	*connect.Response[pb.ExportCharacterResponse], error) {
	resp, err := c.character.ExportCharacter(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ExportCharacter", req.Msg, resp, err)
	}
	return doResponse(ctx, "ExportCharacter", req.Msg, resp)
}

// RestoreCharacterWorldInfo 恢复角色的世界信息
func (c *CharacterServer) RestoreCharacterWorldInfo(ctx context.Context, req *connect.Request[pb.RestoreCharacterWorldInfoRequest]) (
	*connect.Response[pb.RestoreCharacterWorldInfoResponse], error) {
	resp, err := c.character.RestoreCharacterWorldInfo(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "RestoreCharacterWorldInfo", req.Msg, resp, err)
	}
	return doResponse(ctx, "RestoreCharacterWorldInfo", req.Msg, resp)
}

// doResponse 统一处理成功响应并记录日志
func doResponse[T any](_ context.Context, interfaceName string, req any, rsp *T) (
	*connect.Response[T], error) {
	log.Infof("接口: %s, 请求体为: %v, 响应体为: %v", interfaceName, req, rsp)
	return connect.NewResponse(rsp), nil
}

// doResponseExp 统一处理错误响应并记录日志
func doResponseExp[T any](_ context.Context, interfaceName string, req any, rsp *T, err error) (
	*connect.Response[T], error) {
	log.Errorf("接口: %s, 请求体为: %v, 响应体为: %v, 错误为: %v", interfaceName, req, rsp, err)
	return connect.NewResponse(rsp), err
}
