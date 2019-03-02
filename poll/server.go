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

	s.logger.Info().
		Str("pollId", poll.Id).
		Str("title", poll.Title).
		Str("creatorId", poll.CreatorId).
		Int32("allowedUniqueVotes", poll.AllowedUniqueVotes).
		Msg("poll created")

	return &pollpb.CreatePollReply{Poll: poll}, nil
}

func (s Server) Poll(ctx context.Context, req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	poll, err := s.repository.Poll(req.Id)
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
	creator, err := s.repository.PollCreator(req.Id)
	if err != nil {
		return nil, err
	}

	if req.UserId != creator {
		reply := pollpb.EndPollReply{
			Success: false,
			Reason:  "only the poll creator may end the poll",
		}

		return &reply, nil
	}

	poll, err := s.repository.EndPoll(req.Id)
	if err != nil {
		return nil, err
	}

	s.logger.Info().
		Str("pollId", poll.Id).
		Str("user", req.UserId).
		Msg("poll has ended")

	return &pollpb.EndPollReply{Poll: poll, Success: true}, nil
}