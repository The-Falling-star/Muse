package character

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type Character interface {
	ListCharacters(ctx context.Context, req *pb.ListCharactersRequest) (*pb.ListCharactersResponse, error)
	GetCharacter(ctx context.Context, req *pb.GetCharacterRequest) (*pb.GetCharacterResponse, error)
	CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error)
	UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterRequest) (*pb.UpdateCharacterResponse, error)
	DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterRequest) (*pb.DeleteCharacterResponse, error)
	ImportCharacter(ctx context.Context, req *pb.ImportCharacterRequest) (*pb.ImportCharacterResponse, error)
	ExportCharacter(ctx context.Context, req *pb.ExportCharacterRequest) (*pb.ExportCharacterResponse, error)
	RestoreCharacterWorldInfo(ctx context.Context, req *pb.RestoreCharacterWorldInfoRequest) (*pb.RestoreCharacterWorldInfoResponse, error)
}

func NewCharacter() Character {
	return newCharacter()
}
