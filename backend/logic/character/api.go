package character

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// Character 定义了角色服务的接口
type Character interface {
	// ListCharacters 获取角色列表
	ListCharacters(ctx context.Context, req *pb.ListCharactersReq) (*pb.ListCharactersRsp, error)
	// GetCharacter 获取指定角色详情
	GetCharacter(ctx context.Context, req *pb.GetCharacterReq) (*pb.GetCharacterRsp, error)
	// CreateCharacter 创建新角色
	CreateCharacter(ctx context.Context, req *pb.CreateCharacterReq) (*pb.CreateCharacterRsp, error)
	// UpdateCharacter 更新指定角色
	UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterReq) (*pb.UpdateCharacterRsp, error)
	// DeleteCharacter 删除指定角色
	DeleteCharacter(ctx context.Context, req *pb.DeleteCharacterReq) (*pb.DeleteCharacterRsp, error)
	// ImportCharacter 导入角色
	ImportCharacter(ctx context.Context, req *pb.ImportCharacterReq) (*pb.ImportCharacterRsp, error)
	// ExportCharacter 导出角色
	ExportCharacter(ctx context.Context, req *pb.ExportCharacterReq) (*pb.ExportCharacterRsp, error)
	// RestoreCharacterWorldInfo 恢复角色的世界信息
	RestoreCharacterWorldInfo(ctx context.Context, req *pb.RestoreCharacterWorldInfoReq) (*pb.RestoreCharacterWorldInfoRsp, error)
}

// NewCharacter 创建一个新的Character实例
func NewCharacter() Character {
	return newCharacter()
}
