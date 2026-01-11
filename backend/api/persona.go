package api

import (
	"context"

	"connectrpc.com/connect"
	pb "github.com/ling/muse/gen/muse"
	"github.com/ling/muse/logic/persona"
)

type PersonaServer struct {
	persona persona.Persona
}

func NewPersonaServer() *PersonaServer {
	return &PersonaServer{
		persona: persona.NewPersona(),
	}
}

func (p *PersonaServer) ListPersonas(ctx context.Context, req *connect.Request[pb.ListPersonasRequest]) (
	*connect.Response[pb.ListPersonasResponse], error) {
	resp, err := p.persona.ListPersonas(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PersonaServer) GetPersona(ctx context.Context, req *connect.Request[pb.GetPersonaRequest]) (
	*connect.Response[pb.GetPersonaResponse], error) {
	resp, err := p.persona.GetPersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PersonaServer) CreatePersona(ctx context.Context, req *connect.Request[pb.CreatePersonaRequest]) (
	*connect.Response[pb.CreatePersonaResponse], error) {
	resp, err := p.persona.CreatePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PersonaServer) UpdatePersona(ctx context.Context, req *connect.Request[pb.UpdatePersonaRequest]) (
	*connect.Response[pb.UpdatePersonaResponse], error) {
	resp, err := p.persona.UpdatePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PersonaServer) DeletePersona(ctx context.Context, req *connect.Request[pb.DeletePersonaRequest]) (
	*connect.Response[pb.DeletePersonaResponse], error) {
	resp, err := p.persona.DeletePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (p *PersonaServer) SetActivePersona(ctx context.Context, req *connect.Request[pb.SetActivePersonaRequest]) (
	*connect.Response[pb.SetActivePersonaResponse], error) {
	resp, err := p.persona.SetActivePersona(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
