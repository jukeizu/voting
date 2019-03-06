package poll

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type PollHandler interface {
	Create(*pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error)
	Poll(*pollpb.PollRequest) (*pollpb.PollReply, error)
	End(*pollpb.EndPollRequest) (*pollpb.EndPollReply, error)
}

type Validator interface {
	CanEnd(userId, pollId string) (bool, string)
}

type pollHandler struct {
	logger     zerolog.Logger
	repository Repository
	validator  Validator
}

func NewPollHandler(logger zerolog.Logger, repository Repository, validator Validator) PollHandler {
	return &pollHandler{logger, repository, validator}
}

func (h pollHandler) Create(req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
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

func (h pollHandler) Poll(req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	poll, err := h.repository.Poll(req.Id)
	if err != nil {
		return nil, err
	}

	return &pollpb.PollReply{Poll: poll}, nil
}

func (h pollHandler) End(req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
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
