package poll

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
)

type Server struct {
	pollHandler    PollHandler
	optionsHandler OptionsHandler
}

func NewServer(pollHandler PollHandler, optionsHandler OptionsHandler) Server {
	return Server{pollHandler, optionsHandler}
}

func (s Server) Create(ctx context.Context, req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
	return s.pollHandler.Create(req)
}

func (s Server) Poll(ctx context.Context, req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	return s.pollHandler.Poll(req)
}

func (s Server) Options(ctx context.Context, req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	return s.optionsHandler.Options(req)
}

func (s Server) End(ctx context.Context, req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
	return s.pollHandler.End(req)
}
