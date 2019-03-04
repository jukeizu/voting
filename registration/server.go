package registration

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
)

type Server struct {
	registrationHandler RegistrationHandler
}

func NewServer(registrationHandler RegistrationHandler) Server {
	return Server{registrationHandler}
}

func (s Server) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	return s.registrationHandler.RegisterVoter(req)
}
