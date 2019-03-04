package poll

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type OptionsHandler interface {
	Options(*pollpb.OptionsRequest) (*pollpb.OptionsReply, error)
}

type optionsHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewOptionsHandler(logger zerolog.Logger, repository Repository) OptionsHandler {
	return &optionsHandler{logger, repository}
}

func (h optionsHandler) Options(req *pollpb.OptionsRequest) (*pollpb.OptionsReply, error) {
	options, err := h.repository.Options(req.PollId)
	if err != nil {
		return nil, err
	}

	return &pollpb.OptionsReply{Options: options}, nil
}
