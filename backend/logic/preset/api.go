package preset

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

// Preset 定义了预设服务的接口（包含预设和提示项的管理）
type Preset interface {
	// ListPresets 获取预设列表
	ListPresets(ctx context.Context, req *pb.ListPresetsRequest) (*pb.ListPresetsResponse, error)
	// GetPreset 获取指定预设详情
	GetPreset(ctx context.Context, req *pb.GetPresetRequest) (*pb.GetPresetResponse, error)
	// CreatePreset 创建新预设
	CreatePreset(ctx context.Context, req *pb.CreatePresetRequest) (*pb.CreatePresetResponse, error)
	// UpdatePreset 更新指定预设
	UpdatePreset(ctx context.Context, req *pb.UpdatePresetRequest) (*pb.UpdatePresetResponse, error)
	// DeletePreset 删除指定预设
	DeletePreset(ctx context.Context, req *pb.DeletePresetRequest) (*pb.DeletePresetResponse, error)
	// SetActivePreset 设置当前激活的预设
	SetActivePreset(ctx context.Context, req *pb.SetActivePresetRequest) (*pb.SetActivePresetResponse, error)
	// ImportPreset 导入预设
	ImportPreset(ctx context.Context, req *pb.ImportPresetRequest) (*pb.ImportPresetResponse, error)
	// ExportPreset 导出预设
	ExportPreset(ctx context.Context, req *pb.ExportPresetRequest) (*pb.ExportPresetResponse, error)

	// ListPromptItems 获取提示项列表
	ListPromptItems(ctx context.Context, req *pb.ListPromptItemsRequest) (*pb.ListPromptItemsResponse, error)
	// AddPromptItem 添加新的提示项
	AddPromptItem(ctx context.Context, req *pb.AddPromptItemRequest) (*pb.AddPromptItemResponse, error)
	// UpdatePromptItem 更新指定提示项
	UpdatePromptItem(ctx context.Context, req *pb.UpdatePromptItemRequest) (*pb.UpdatePromptItemResponse, error)
	// DeletePromptItem 删除指定提示项
	DeletePromptItem(ctx context.Context, req *pb.DeletePromptItemRequest) (*pb.DeletePromptItemResponse, error)
	// UpdatePromptItemsOrder 更新提示项的排序
	UpdatePromptItemsOrder(ctx context.Context, req *pb.UpdatePromptItemsOrderRequest) (*pb.UpdatePromptItemsOrderResponse, error)
}

// NewPreset 创建一个新的Preset实例
func NewPreset() Preset {
	return newPreset()
}
