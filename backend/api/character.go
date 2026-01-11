package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/character"
)

type CharacterServer struct {
	character character.Character
}

func NewCharacterServer() *CharacterServer {
	return &CharacterServer{
		character: character.NewCharacter(),
	}
}

func (c *CharacterServer) ListCharacters(ctx context.Context, req *connect.Request[pb.ListCharactersRequest]) (
	*connect.Response[pb.ListCharactersResponse], error) {
	resp, err := c.character.ListCharacters(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *CharacterServer) GetCharacter(ctx context.Context, req *connect.Request[pb.GetCharacterRequest]) (
	*connect.Response[pb.GetCharacterResponse], error) {
	resp, err := c.character.GetCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *CharacterServer) CreateCharacter(ctx context.Context, req *connect.Request[pb.CreateCharacterRequest]) (
	*connect.Response[pb.CreateCharacterResponse], error) {
	resp, err := c.character.CreateCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *CharacterServer) UpdateCharacter(ctx context.Context, req *connect.Request[pb.UpdateCharacterRequest]) (
	*connect.Response[pb.UpdateCharacterResponse], error) {
	resp, err := c.character.UpdateCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *CharacterServer) DeleteCharacter(ctx context.Context, req *connect.Request[pb.DeleteCharacterRequest]) (
	*connect.Response[pb.DeleteCharacterResponse], error) {
	resp, err := c.character.DeleteCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *CharacterServer) ImportCharacter(ctx context.Context, req *connect.Request[pb.ImportCharacterRequest]) (
	*connect.Response[pb.ImportCharacterResponse], error) {
	resp, err := c.character.ImportCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *CharacterServer) ExportCharacter(ctx context.Context, req *connect.Request[pb.ExportCharacterRequest]) (
	*connect.Response[pb.ExportCharacterResponse], error) {
	resp, err := c.character.ExportCharacter(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (c *CharacterServer) RestoreCharacterWorldInfo(ctx context.Context, req *connect.Request[pb.RestoreCharacterWorldInfoRequest]) (
	*connect.Response[pb.RestoreCharacterWorldInfoResponse], error) {
	resp, err := c.character.RestoreCharacterWorldInfo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
