package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/preset"
)

// PresetServer 预设服务（包含预设和提示项的管理）
type PresetServer struct {
	preset preset.Preset
}

// NewPresetServer 创建一个新的PresetServer实例
func NewPresetServer() *PresetServer {
	return &PresetServer{
		preset: preset.NewPreset(),
	}
}

// ListPresets 获取预设列表
func (p *PresetServer) ListPresets(ctx context.Context, req *connect.Request[pb.ListPresetsRequest]) (
	*connect.Response[pb.ListPresetsResponse], error) {
	resp, err := p.preset.ListPresets(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ListPresets", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListPresets", req.Msg, resp)
}

// GetPreset 获取指定预设详情
func (p *PresetServer) GetPreset(ctx context.Context, req *connect.Request[pb.GetPresetRequest]) (
	*connect.Response[pb.GetPresetResponse], error) {
	resp, err := p.preset.GetPreset(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "GetPreset", req.Msg, resp, err)
	}
	return doResponse(ctx, "GetPreset", req.Msg, resp)
}

// CreatePreset 创建新预设
func (p *PresetServer) CreatePreset(ctx context.Context, req *connect.Request[pb.CreatePresetRequest]) (
	*connect.Response[pb.CreatePresetResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "CreatePreset", req.Msg, (*pb.CreatePresetResponse)(nil), err)
	}

	resp, err := p.preset.CreatePreset(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "CreatePreset", req.Msg, resp, err)
	}
	return doResponse(ctx, "CreatePreset", req.Msg, resp)
}

// UpdatePreset 更新指定预设
func (p *PresetServer) UpdatePreset(ctx context.Context, req *connect.Request[pb.UpdatePresetRequest]) (
	*connect.Response[pb.UpdatePresetResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdatePreset", req.Msg, (*pb.UpdatePresetResponse)(nil), err)
	}

	resp, err := p.preset.UpdatePreset(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdatePreset", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdatePreset", req.Msg, resp)
}

// DeletePreset 删除指定预设
func (p *PresetServer) DeletePreset(ctx context.Context, req *connect.Request[pb.DeletePresetRequest]) (
	*connect.Response[pb.DeletePresetResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "DeletePreset", req.Msg, (*pb.DeletePresetResponse)(nil), err)
	}

	resp, err := p.preset.DeletePreset(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeletePreset", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeletePreset", req.Msg, resp)
}

// SetActivePreset 设置当前激活的预设
func (p *PresetServer) SetActivePreset(ctx context.Context, req *connect.Request[pb.SetActivePresetRequest]) (
	*connect.Response[pb.SetActivePresetResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "SetActivePreset", req.Msg, (*pb.SetActivePresetResponse)(nil), err)
	}

	resp, err := p.preset.SetActivePreset(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "SetActivePreset", req.Msg, resp, err)
	}
	return doResponse(ctx, "SetActivePreset", req.Msg, resp)
}

// ImportPreset 导入预设
func (p *PresetServer) ImportPreset(ctx context.Context, req *connect.Request[pb.ImportPresetRequest]) (
	*connect.Response[pb.ImportPresetResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "ImportPreset", req.Msg, (*pb.ImportPresetResponse)(nil), err)
	}

	resp, err := p.preset.ImportPreset(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ImportPreset", req.Msg, resp, err)
	}
	return doResponse(ctx, "ImportPreset", req.Msg, resp)
}

// ExportPreset 导出预设
func (p *PresetServer) ExportPreset(ctx context.Context, req *connect.Request[pb.ExportPresetRequest]) (
	*connect.Response[pb.ExportPresetResponse], error) {
	resp, err := p.preset.ExportPreset(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ExportPreset", req.Msg, resp, err)
	}
	return doResponse(ctx, "ExportPreset", req.Msg, resp)
}

// ListPromptItems 获取提示项列表
func (p *PresetServer) ListPromptItems(ctx context.Context, req *connect.Request[pb.ListPromptItemsRequest]) (
	*connect.Response[pb.ListPromptItemsResponse], error) {
	resp, err := p.preset.ListPromptItems(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "ListPromptItems", req.Msg, resp, err)
	}
	return doResponse(ctx, "ListPromptItems", req.Msg, resp)
}

// AddPromptItem 添加新的提示项
func (p *PresetServer) AddPromptItem(ctx context.Context, req *connect.Request[pb.AddPromptItemRequest]) (
	*connect.Response[pb.AddPromptItemResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "AddPromptItem", req.Msg, (*pb.AddPromptItemResponse)(nil), err)
	}

	resp, err := p.preset.AddPromptItem(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "AddPromptItem", req.Msg, resp, err)
	}
	return doResponse(ctx, "AddPromptItem", req.Msg, resp)
}

// UpdatePromptItem 更新指定提示项
func (p *PresetServer) UpdatePromptItem(ctx context.Context, req *connect.Request[pb.UpdatePromptItemRequest]) (
	*connect.Response[pb.UpdatePromptItemResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdatePromptItem", req.Msg, (*pb.UpdatePromptItemResponse)(nil), err)
	}

	resp, err := p.preset.UpdatePromptItem(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdatePromptItem", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdatePromptItem", req.Msg, resp)
}

// DeletePromptItem 删除指定提示项
func (p *PresetServer) DeletePromptItem(ctx context.Context, req *connect.Request[pb.DeletePromptItemRequest]) (
	*connect.Response[pb.DeletePromptItemResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "DeletePromptItem", req.Msg, (*pb.DeletePromptItemResponse)(nil), err)
	}

	resp, err := p.preset.DeletePromptItem(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "DeletePromptItem", req.Msg, resp, err)
	}
	return doResponse(ctx, "DeletePromptItem", req.Msg, resp)
}

// UpdatePromptItemsOrder 更新提示项的排序
func (p *PresetServer) UpdatePromptItemsOrder(ctx context.Context,
	req *connect.Request[pb.UpdatePromptItemsOrderRequest]) (
	*connect.Response[pb.UpdatePromptItemsOrderResponse], error) {
	// 开启事务
	var err error
	ctx, err = beginTransaction(ctx)
	if err != nil {
		return doResponseExp(ctx, "UpdatePromptItemsOrder", req.Msg, (*pb.UpdatePromptItemsOrderResponse)(nil), err)
	}

	resp, err := p.preset.UpdatePromptItemsOrder(ctx, req.Msg)
	if err != nil {
		return doResponseExp(ctx, "UpdatePromptItemsOrder", req.Msg, resp, err)
	}
	return doResponse(ctx, "UpdatePromptItemsOrder", req.Msg, resp)
}
