package preset

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type presetImpl struct {
}

func newPreset() *presetImpl {
	return &presetImpl{}
}

func (p *presetImpl) ListPresets(ctx context.Context, req *pb.ListPresetsRequest) (*pb.ListPresetsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) GetPreset(ctx context.Context, req *pb.GetPresetRequest) (*pb.GetPresetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) CreatePreset(ctx context.Context, req *pb.CreatePresetRequest) (*pb.CreatePresetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) UpdatePreset(ctx context.Context, req *pb.UpdatePresetRequest) (*pb.UpdatePresetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) DeletePreset(ctx context.Context, req *pb.DeletePresetRequest) (*pb.DeletePresetResponse, error) {
	//TODO implement me
	panic("implement me")
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

func (p *presetImpl) ListPromptItems(ctx context.Context, req *pb.ListPromptItemsRequest) (*pb.ListPromptItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) AddPromptItem(ctx context.Context, req *pb.AddPromptItemRequest) (*pb.AddPromptItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) UpdatePromptItem(ctx context.Context, req *pb.UpdatePromptItemRequest) (*pb.UpdatePromptItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) DeletePromptItem(ctx context.Context, req *pb.DeletePromptItemRequest) (*pb.DeletePromptItemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *presetImpl) UpdatePromptItemsOrder(ctx context.Context, req *pb.UpdatePromptItemsOrderRequest) (*pb.UpdatePromptItemsOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}
