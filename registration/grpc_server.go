package registration

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
)

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) RegisterVoter(ctx context.Context, req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	return s.service.RegisterVoter(req)
}
