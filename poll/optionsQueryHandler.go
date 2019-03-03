package poll

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type OptionsQueryHandler interface {
	Handle(*pollpb.OptionsRequest) (*pollpb.OptionsReply, error)
}

type optionsQueryHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewOptionsQueryHandler(logger zerolog.Logger, repository Repository) OptionsQueryHandler {
	return &optionsQueryHandler{logger, repository}
}

func (h optionsQueryHandler) Handle(req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	options, err := h.repository.Options(req.PollId)
	if err != nil {
		return nil, err
	}

	return &pollpb.OptionsReply{Options: options}, nil
}
