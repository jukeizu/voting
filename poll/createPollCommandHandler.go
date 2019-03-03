package poll

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type CreatePollCommandHandler interface {
	Handle(*pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error)
}

type createPollCommandHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewCreatePollCommandHandler(logger zerolog.Logger, repository Repository) CreatePollCommandHandler {
	return &createPollCommandHandler{logger, repository}
}

func (h createPollCommandHandler) Handle(req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
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
