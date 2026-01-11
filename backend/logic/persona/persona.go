package persona

import (
	"context"

	pb "github.com/ling/muse/gen/muse"
)

type personaImpl struct {
}

func newPersona() *personaImpl {
	return &personaImpl{}
}

func (p *personaImpl) ListPersonas(ctx context.Context, req *pb.ListPersonasRequest) (*pb.ListPersonasResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *personaImpl) GetPersona(ctx context.Context, req *pb.GetPersonaRequest) (*pb.GetPersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *personaImpl) CreatePersona(ctx context.Context, req *pb.CreatePersonaRequest) (*pb.CreatePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *personaImpl) UpdatePersona(ctx context.Context, req *pb.UpdatePersonaRequest) (*pb.UpdatePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *personaImpl) DeletePersona(ctx context.Context, req *pb.DeletePersonaRequest) (*pb.DeletePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *personaImpl) SetActivePersona(ctx context.Context, req *pb.SetActivePersonaRequest) (*pb.SetActivePersonaResponse, error) {
	//TODO implement me
	panic("implement me")
}
