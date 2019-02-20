package poll

import (
	"context"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type Server struct {
	logger     zerolog.Logger
	repository Repository
}

func NewServer(logger zerolog.Logger, repository Repository) Server {
	return Server{logger, repository}
}

func (s Server) Create(ctx context.Context, req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
	poll, err := s.repository.CreatePoll(req)
	if err != nil {
		return nil, err
	}

	return &pollpb.CreatePollReply{Poll: poll}, nil
}

func (s Server) Poll(ctx context.Context, req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	poll, err := s.repository.Poll(req)
	if err != nil {
		return nil, err
	}

	return &pollpb.PollReply{Poll: poll}, nil
}

func (s Server) Options(ctx context.Context, req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	options, err := s.repository.Options(req.PollId)
	if err != nil {
		return nil, err
	}

	return &pollpb.OptionsReply{Options: options}, nil
}

func (s Server) End(ctx context.Context, req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
	return nil, nil
}
