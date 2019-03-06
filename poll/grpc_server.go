package poll

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
)

type GrpcServer struct {
	service Service
}

func NewGrpcServer(service Service) GrpcServer {
	return GrpcServer{service}
}

func (s GrpcServer) Create(ctx context.Context, req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
	return s.service.Create(req)
}

func (s GrpcServer) Poll(ctx context.Context, req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	return s.service.Poll(req)
}

func (s GrpcServer) Options(ctx context.Context, req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	return s.service.Options(req)
}

func (s GrpcServer) End(ctx context.Context, req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
	return s.service.End(req)
}
