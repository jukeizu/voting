package registration

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/rs/zerolog"
)

type Server struct {
	logger   zerolog.Logger
	mediator Mediator
}

func NewServer(logger zerolog.Logger, mediator Mediator) Server {
	return Server{logger, mediator}
}

func (s Server) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
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
