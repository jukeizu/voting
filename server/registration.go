package server

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
)

type RegistrationServer struct {
	mediator Mediator
}

func NewRegistrationServer(mediator Mediator) RegistrationServer {
	return RegistrationServer{mediator}
}

func (s RegistrationServer) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	resp, err := s.mediator.Send(req)
	if err != nil {
		return nil, err
	}

	reply, ok := resp.(*registrationpb.RegisterVoterReply)
	if !ok {
		return nil, nil
	}

	return reply, nil
}
