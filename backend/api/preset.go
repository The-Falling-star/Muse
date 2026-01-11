package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/preset"
)

type PresetServer struct {
	preset preset.Preset
}

func NewPresetServer() *PresetServer {
	return &PresetServer{
		preset: preset.NewPreset(),
	}
}

func (p *PresetServer) ListPresets(ctx context.Context, req *connect.Request[pb.ListPresetsRequest]) (
	*connect.Response[pb.ListPresetsResponse], error) {
	resp, err := p.preset.ListPresets(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PresetServer) GetPreset(ctx context.Context, req *connect.Request[pb.GetPresetRequest]) (
	*connect.Response[pb.GetPresetResponse], error) {
	resp, err := p.preset.GetPreset(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PresetServer) CreatePreset(ctx context.Context, req *connect.Request[pb.CreatePresetRequest]) (
	*connect.Response[pb.CreatePresetResponse], error) {
	resp, err := p.preset.CreatePreset(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PresetServer) UpdatePreset(ctx context.Context, req *connect.Request[pb.UpdatePresetRequest]) (
	*connect.Response[pb.UpdatePresetResponse], error) {
	resp, err := p.preset.UpdatePreset(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PresetServer) DeletePreset(ctx context.Context, req *connect.Request[pb.DeletePresetRequest]) (
	*connect.Response[pb.DeletePresetResponse], error) {
	resp, err := p.preset.DeletePreset(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PresetServer) SetActivePreset(ctx context.Context, req *connect.Request[pb.SetActivePresetRequest]) (
	*connect.Response[pb.SetActivePresetResponse], error) {
	resp, err := p.preset.SetActivePreset(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PresetServer) ImportPreset(ctx context.Context, req *connect.Request[pb.ImportPresetRequest]) (
	*connect.Response[pb.ImportPresetResponse], error) {
	resp, err := p.preset.ImportPreset(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PresetServer) ExportPreset(ctx context.Context, req *connect.Request[pb.ExportPresetRequest]) (
	*connect.Response[pb.ExportPresetResponse], error) {
	resp, err := p.preset.ExportPreset(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
