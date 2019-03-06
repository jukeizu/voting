package poll

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
)

type GrpcServer struct {
	pollHandler PollHandler
}

func NewGrpcServer(pollHandler PollHandler) GrpcServer {
	return GrpcServer{pollHandler}
}

func (s GrpcServer) Create(ctx context.Context, req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
	return s.pollHandler.Create(req)
}

func (s GrpcServer) Poll(ctx context.Context, req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	return s.pollHandler.Poll(req)
}

func (s GrpcServer) Options(ctx context.Context, req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	return s.pollHandler.Options(req)
}

func (s GrpcServer) End(ctx context.Context, req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
	return s.pollHandler.End(req)
}
