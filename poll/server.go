package poll

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
)

type Server struct {
	createPollCommandHandler CreatePollCommandHandler
	pollQueryHandler         PollQueryHandler
	optionsQueryHandler      OptionsQueryHandler
	endPollCommandHandler    EndPollCommandHandler
}

func NewServer(
	createPollCommandHandler CreatePollCommandHandler,
	pollQueryHandler PollQueryHandler,
	optionsQueryHandler OptionsQueryHandler,
	endPollCommandHandler EndPollCommandHandler,
) Server {
	return Server{
		createPollCommandHandler,
		pollQueryHandler,
		optionsQueryHandler,
		endPollCommandHandler,
	}
}

func (s Server) Create(ctx context.Context, req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
	return s.createPollCommandHandler.Handle(req)
}

func (s Server) Poll(ctx context.Context, req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	return s.pollQueryHandler.Handle(req)
}

func (s Server) Options(ctx context.Context, req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	return s.optionsQueryHandler.Handle(req)
}

func (s Server) End(ctx context.Context, req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
	return s.endPollCommandHandler.Handle(req)
}
