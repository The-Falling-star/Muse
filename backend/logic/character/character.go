package character

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type characterImpl struct {
}

func newCharacter() *characterImpl {
	return &characterImpl{}
}

func (c *characterImpl) ListCharacters(ctx context.Context, req *pb.ListCharactersRequest) (*pb.ListCharactersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) GetCharacter(ctx context.Context, req *pb.GetCharacterRequest) (*pb.GetCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterRequest) (*pb.UpdateCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterRequest) (*pb.DeleteCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) ImportCharacter(ctx context.Context, req *pb.ImportCharacterRequest) (*pb.ImportCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) ExportCharacter(ctx context.Context, req *pb.ExportCharacterRequest) (*pb.ExportCharacterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *characterImpl) RestoreCharacterWorldInfo(ctx context.Context, req *pb.RestoreCharacterWorldInfoRequest) (*pb.RestoreCharacterWorldInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}
