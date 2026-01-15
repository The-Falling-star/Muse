package preset

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

type presetImpl struct {
	presetRepo *database.PresetRepo
}

func newPreset() *presetImpl {
	return &presetImpl{
		presetRepo: &database.PresetRepo{},
	}
}

func (p *presetImpl) ListPresets(ctx context.Context, req *pb.ListPresetsRequest) (*pb.ListPresetsResponse, error) {
	// 从数据库获取预设列表（不分页）
	presets, _, err := p.presetRepo.List(defaultUserID, 1, 1000)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 转换为pb格式
	pbPresets := make([]*pb.Preset, 0, len(presets))
	for _, preset := range presets {
		pbPresets = append(pbPresets, convert.PresetEntityToPb(preset))
	}

	return &pb.ListPresetsResponse{
		Presets: pbPresets,
	}, nil
}

func (p *presetImpl) GetPreset(ctx context.Context, req *pb.GetPresetRequest) (*pb.GetPresetResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	// 从数据库获取预设
	preset, err := p.presetRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if preset == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("预设不存在"))
	}

	return &pb.GetPresetResponse{
		Preset: convert.PresetEntityToPb(preset),
	}, nil
}

func (p *presetImpl) CreatePreset(ctx context.Context, req *pb.CreatePresetRequest) (*pb.CreatePresetResponse, error) {
	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("预设名称不能为空"))
	}

	// 构建预设实体
	preset := &entity.Preset{
		UserID:           defaultUserID,
		Name:             name,
		Temperature:      req.GetTemperature(),
		TopP:             req.GetTopP(),
		TopK:             int(req.GetTopK()),
		MaxTokens:        int(req.GetMaxTokens()),
		FrequencyPenalty: req.GetFrequencyPenalty(),
		PresencePenalty:  req.GetPresencePenalty(),
		Version:          1,
	}

	// 保存到数据库
	if err := p.presetRepo.Create(preset); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取完整数据（包含关联）
	fullPreset, err := p.presetRepo.GetByID(preset.ID, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.CreatePresetResponse{
		Preset: convert.PresetEntityToPb(fullPreset),
	}, nil
}

func (p *presetImpl) UpdatePreset(ctx context.Context, req *pb.UpdatePresetRequest) (*pb.UpdatePresetResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	// 参数校验
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("预设名称不能为空"))
	}

	// 获取当前预设
	preset, err := p.presetRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if preset == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("预设不存在"))
	}

	// 检查版本号
	if int64(preset.Version) != req.GetVersion() {
		return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("数据已被修改，请刷新后重试"))
	}

	// 更新预设字段
	preset.Name = name
	preset.Temperature = req.Temperature
	preset.TopP = req.TopP
	preset.TopK = int(req.TopK)
	preset.MaxTokens = int(req.MaxTokens)
	preset.FrequencyPenalty = req.FrequencyPenalty
	preset.PresencePenalty = req.PresencePenalty

	// 更新数据库
	if err := p.presetRepo.Update(preset); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的预设（包含关联数据）
	updatedPreset, err := p.presetRepo.GetByID(id, defaultUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdatePresetResponse{
		Preset: convert.PresetEntityToPb(updatedPreset),
	}, nil
}

func (p *presetImpl) DeletePreset(ctx context.Context, req *pb.DeletePresetRequest) (*pb.DeletePresetResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	// 删除预设（会级联删除关联的提示项）
	if err := p.presetRepo.Delete(id, defaultUserID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.DeletePresetResponse{}, nil
}

func (p *presetImpl) ListPromptItems(ctx context.Context, req *pb.ListPromptItemsRequest) (*pb.ListPromptItemsResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	// 从数据库获取提示项列表
	items, err := p.presetRepo.ListPromptItems(presetID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 转换为pb格式
	pbItems := make([]*pb.PromptItem, 0, len(items))
	for _, item := range items {
		pbItems = append(pbItems, convert.PromptItemEntityToPb(item))
	}

	return &pb.ListPromptItemsResponse{
		Items: pbItems,
	}, nil
}

func (p *presetImpl) AddPromptItem(ctx context.Context, req *pb.AddPromptItemRequest) (*pb.AddPromptItemResponse, error) {
	// 参数校验
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("提示项名称不能为空"))
	}

	// 构建提示项实体
	item := &entity.PromptItem{
		PresetID:          presetID,
		Identifier:        strings.TrimSpace(req.Identifier),
		Name:              name,
		Content:           req.GetContent(),
		Role:              req.GetRole(),
		IsEnabled:         req.GetIsEnabled(),
		InjectionPosition: req.GetInjectionPosition(),
		InjectionDepth:    int(req.GetInjectionDepth()),
		ForbidOverrides:   req.GetForbidOverrides(),
		SortOrder:         int(req.GetSortOrder()),
	}

	// 保存到数据库
	if err := p.presetRepo.CreatePromptItem(item); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取完整数据
	fullItem, err := p.presetRepo.GetPromptItemByID(item.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.AddPromptItemResponse{
		Item: convert.PromptItemEntityToPb(fullItem),
	}, nil
}

func (p *presetImpl) UpdatePromptItem(ctx context.Context, req *pb.UpdatePromptItemRequest) (*pb.UpdatePromptItemResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的提示项ID"))
	}

	// 参数校验
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("提示项名称不能为空"))
	}

	// 获取当前提示项
	item, err := p.presetRepo.GetPromptItemByID(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if item == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("提示项不存在"))
	}

	// 更新提示项字段
	item.Identifier = strings.TrimSpace(req.Identifier)
	item.Name = name
	if req.Content != nil {
		item.Content = *req.Content
	}
	item.Role = req.Role
	item.IsEnabled = req.IsEnabled
	item.InjectionPosition = req.InjectionPosition
	item.InjectionDepth = int(req.InjectionDepth)
	item.ForbidOverrides = req.ForbidOverrides
	item.SortOrder = int(req.SortOrder)

	// 更新数据库
	if err := p.presetRepo.UpdatePromptItem(item); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 重新获取更新后的提示项
	updatedItem, err := p.presetRepo.GetPromptItemByID(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdatePromptItemResponse{
		Item: convert.PromptItemEntityToPb(updatedItem),
	}, nil
}

func (p *presetImpl) DeletePromptItem(ctx context.Context, req *pb.DeletePromptItemRequest) (*pb.DeletePromptItemResponse, error) {
	id := int(req.GetId())
	if id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的提示项ID"))
	}

	// 删除提示项
	if err := p.presetRepo.DeletePromptItem(id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.DeletePromptItemResponse{}, nil
}

func (p *presetImpl) UpdatePromptItemsOrder(ctx context.Context, req *pb.UpdatePromptItemsOrderRequest) (*pb.UpdatePromptItemsOrderResponse, error) {
	presetID := int(req.GetPresetId())
	if presetID <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("无效的预设ID"))
	}

	if len(req.ItemIds) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("排序数据不能为空"))
	}

	// 构建排序映射
	itemOrders := make(map[int]int)
	for i, itemID := range req.ItemIds {
		itemOrders[int(itemID)] = i
	}

	// 更新排序
	if err := p.presetRepo.UpdatePromptItemsOrder(presetID, itemOrders); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &pb.UpdatePromptItemsOrderResponse{}, nil
}

func (p *presetImpl) SetActivePreset(ctx context.Context, req *pb.SetActivePresetRequest) (*pb.SetActivePresetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) ImportPreset(ctx context.Context, req *pb.ImportPresetRequest) (*pb.ImportPresetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) ExportPreset(ctx context.Context, req *pb.ExportPresetRequest) (*pb.ExportPresetResponse, error) {
	//TODO implement me
	panic("implement me")
}
