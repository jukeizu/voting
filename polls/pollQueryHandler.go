package polls

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type PollQueryHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewPollQueryHandler(logger zerolog.Logger, repository Repository) PollQueryHandler {
	return PollQueryHandler{logger, repository}
}

func (h PollQueryHandler) Handle(request interface{}) (interface{}, error) {
	req, ok := request.(*pollpb.PollRequest)
	if !ok {
		return nil, nil
	}

	poll, err := h.repository.Poll(req.Id)
	if err != nil {
		return nil, err
	}

	pbPoll := mapToPbPoll(poll)

	return &pollpb.PollReply{Poll: pbPoll}, nil
}
