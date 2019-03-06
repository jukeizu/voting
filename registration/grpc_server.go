package registration

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
)

type GrpcServer struct {
	registrationHandler RegistrationHandler
}

func NewGrpcServer(registrationHandler RegistrationHandler) GrpcServer {
	return GrpcServer{registrationHandler}
}

func (s GrpcServer) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	return s.registrationHandler.RegisterVoter(req)
}
