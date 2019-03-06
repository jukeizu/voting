package poll

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type Service interface {
	Create(*pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error)
	Poll(*pollpb.PollRequest) (*pollpb.PollReply, error)
	Options(*pollpb.OptionsRequest) (*pollpb.OptionsReply, error)
	End(*pollpb.EndPollRequest) (*pollpb.EndPollReply, error)
}

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
}

func NewDefaultService(logger zerolog.Logger, repository Repository) Service {
	return &DefaultService{logger, repository}
}

func (h DefaultService) Create(req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
	poll, err := h.repository.CreatePoll(req)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Str("title", poll.Title).
		Str("creatorId", poll.CreatorId).
		Int32("allowedUniqueVotes", poll.AllowedUniqueVotes).
		Msg("poll created")

	return &pollpb.CreatePollReply{Poll: poll}, nil
}

func (h DefaultService) Poll(req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	poll, err := h.repository.Poll(req.Id)
	if err != nil {
		return nil, err
	}

	return &pollpb.PollReply{Poll: poll}, nil
}

func (h DefaultService) Options(req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	options, err := h.repository.Options(req.PollId)
	if err != nil {
		return nil, err
	}

	return &pollpb.OptionsReply{Options: options}, nil
}

func (h DefaultService) End(req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
	creator, err := h.repository.PollCreator(req.Id)
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

	poll, err := h.repository.EndPoll(req.Id)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Str("user", req.UserId).
		Msg("poll has ended")

	return &pollpb.EndPollReply{Poll: poll, Success: true}, nil
}
