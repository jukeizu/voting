package poll

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type EndPollCommandHandler interface {
	Handle(*pollpb.EndPollRequest) (*pollpb.EndPollReply, error)
}

type endPollCommandHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewEndPollCommandHandler(logger zerolog.Logger, repository Repository) EndPollCommandHandler {
	return &endPollCommandHandler{logger, repository}
}

func (h endPollCommandHandler) Handle(req *pollpb.EndPollRequest) (*pollpb.EndPollReply, error) {
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
