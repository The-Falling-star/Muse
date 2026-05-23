package character

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// Character 定义了角色服务的接口
type Character interface {
	// ListCharacters 获取角色列表
	ListCharacters(ctx context.Context, req *pb.ListCharactersRequest) (*pb.ListCharactersResponse, error)
	// GetCharacter 获取指定角色详情
	GetCharacter(ctx context.Context, req *pb.GetCharacterRequest) (*pb.GetCharacterResponse, error)
	// CreateCharacter 创建新角色
	CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error)
	// UpdateCharacter 更新指定角色
	UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterRequest) (*pb.UpdateCharacterResponse, error)
	// DeleteCharacter 删除指定角色
	DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterRequest) (*pb.DeleteCharacterResponse, error)
	// ImportCharacter 导入角色
	ImportCharacter(ctx context.Context, req *pb.ImportCharacterRequest) (*pb.ImportCharacterResponse, error)
	// ExportCharacter 导出角色
	ExportCharacter(ctx context.Context, req *pb.ExportCharacterRequest) (*pb.ExportCharacterResponse, error)
	// RestoreCharacterWorldInfo 恢复角色的世界信息
	RestoreCharacterWorldInfo(ctx context.Context, req *pb.RestoreCharacterWorldInfoRequest) (*pb.RestoreCharacterWorldInfoResponse, error)
}

// NewCharacter 创建一个新的Character实例
func NewCharacter() Character {
	return newCharacter()
}
