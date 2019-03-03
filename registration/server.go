package registration

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
)

type Server struct {
	registerVoterCommandHandler RegisterVoterCommandHandler
}

func NewServer(registerVoterCommandHandler RegisterVoterCommandHandler) Server {
	return Server{registerVoterCommandHandler}
}

func (s Server) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	return s.registerVoterCommandHandler.Handle(req)
}
