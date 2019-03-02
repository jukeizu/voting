package polls

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/rs/zerolog"
)

type OptionsQueryHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewOptionsQueryHandler(logger zerolog.Logger, repository Repository) OptionsQueryHandler {
	return OptionsQueryHandler{logger, repository}
}

func (h OptionsQueryHandler) Handle(request interface{}) (interface{}, error) {
	req, ok := request.(*pollpb.OptionsRequest)
	if !ok {
		return nil, nil
	}

	options, err := h.repository.Options(req.PollId)
	if err != nil {
		return nil, err
	}

	pbOptions := mapToPbOptions(options)

	return &pollpb.OptionsReply{Options: pbOptions}, nil
}
