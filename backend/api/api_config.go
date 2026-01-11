package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/api_config"
)

type APIConfigServer struct {
	apiConfig api_config.APIConfig
}

func NewAPIConfigServer() *APIConfigServer {
	return &APIConfigServer{
		apiConfig: api_config.NewAPIConfig(),
	}
}

func (a *APIConfigServer) ListAPIConfigs(ctx context.Context, req *connect.Request[pb.ListAPIConfigsRequest]) (
	*connect.Response[pb.ListAPIConfigsResponse], error) {
	resp, err := a.apiConfig.ListAPIConfigs(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *APIConfigServer) GetAPIConfig(ctx context.Context, req *connect.Request[pb.GetAPIConfigRequest]) (
	*connect.Response[pb.GetAPIConfigResponse], error) {
	resp, err := a.apiConfig.GetAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *APIConfigServer) CreateAPIConfig(ctx context.Context, req *connect.Request[pb.CreateAPIConfigRequest]) (
	*connect.Response[pb.CreateAPIConfigResponse], error) {
	resp, err := a.apiConfig.CreateAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *APIConfigServer) UpdateAPIConfig(ctx context.Context, req *connect.Request[pb.UpdateAPIConfigRequest]) (
	*connect.Response[pb.UpdateAPIConfigResponse], error) {
	resp, err := a.apiConfig.UpdateAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *APIConfigServer) DeleteAPIConfig(ctx context.Context, req *connect.Request[pb.DeleteAPIConfigRequest]) (
	*connect.Response[pb.DeleteAPIConfigResponse], error) {
	resp, err := a.apiConfig.DeleteAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *APIConfigServer) SetActiveAPIConfig(ctx context.Context, req *connect.Request[pb.SetActiveAPIConfigRequest]) (
	*connect.Response[pb.SetActiveAPIConfigResponse], error) {
	resp, err := a.apiConfig.SetActiveAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (a *APIConfigServer) TestAPIConfig(ctx context.Context, req *connect.Request[pb.TestAPIConfigRequest]) (
	*connect.Response[pb.TestAPIConfigResponse], error) {
	resp, err := a.apiConfig.TestAPIConfig(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
