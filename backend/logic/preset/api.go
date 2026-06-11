package preset

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// Preset 定义了预设服务的接口（包含预设和提示项的管理）
type Preset interface {
	// ListPresets 获取预设列表
	ListPresets(ctx context.Context, req *pb.ListPresetsReq) (*pb.ListPresetsRsp, error)
	// GetPreset 获取指定预设详情
	GetPreset(ctx context.Context, req *pb.GetPresetReq) (*pb.GetPresetRsp, error)
	// CreatePreset 创建新预设
	CreatePreset(ctx context.Context, req *pb.CreatePresetReq) (*pb.CreatePresetRsp, error)
	// UpdatePreset 更新指定预设
	UpdatePreset(ctx context.Context, req *pb.UpdatePresetReq) (*pb.UpdatePresetRsp, error)
	// DeletePreset 删除指定预设
	DeletePreset(ctx context.Context, req *pb.DeletePresetReq) (*pb.DeletePresetRsp, error)
	// SetActivePreset 设置当前激活的预设
	SetActivePreset(ctx context.Context, req *pb.SetActivePresetReq) (*pb.SetActivePresetRsp, error)
	// ImportPreset 导入预设
	ImportPreset(ctx context.Context, req *pb.ImportPresetReq) (*pb.ImportPresetRsp, error)
	// ExportPreset 导出预设
	ExportPreset(ctx context.Context, req *pb.ExportPresetReq) (*pb.ExportPresetRsp, error)

	// ListPromptItems 获取提示项列表
	ListPromptItems(ctx context.Context, req *pb.ListPromptItemsReq) (*pb.ListPromptItemsRsp, error)
	// AddPromptItem 添加新的提示项
	AddPromptItem(ctx context.Context, req *pb.AddPromptItemReq) (*pb.AddPromptItemRsp, error)
	// UpdatePromptItem 更新指定提示项
	UpdatePromptItem(ctx context.Context, req *pb.UpdatePromptItemReq) (*pb.UpdatePromptItemRsp, error)
	// DeletePromptItem 删除指定提示项
	DeletePromptItem(ctx context.Context, req *pb.DeletePromptItemReq) (*pb.DeletePromptItemRsp, error)
	// UpdatePromptItemsOrder 更新提示项的排序
	UpdatePromptItemsOrder(ctx context.Context, req *pb.UpdatePromptItemsOrderReq) (*pb.UpdatePromptItemsOrderRsp, error)
}

// NewPreset 创建一个新的Preset实例
func NewPreset() Preset {
	return newPreset()
}
