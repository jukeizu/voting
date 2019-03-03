package poll

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type PollQueryHandler interface {
	Handle(*pollpb.PollRequest) (*pollpb.PollReply, error)
}

type pollQueryHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewPollQueryHandler(logger zerolog.Logger, repository Repository) PollQueryHandler {
	return pollQueryHandler{logger, repository}
}

func (h pollQueryHandler) Handle(req *pollpb.PollRequest) (*pollpb.PollReply, error) {
	poll, err := h.repository.Poll(req.Id)
	if err != nil {
		return nil, err
	}

	return &pollpb.PollReply{Poll: poll}, nil
}
