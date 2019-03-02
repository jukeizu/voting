package polls

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/domain/entities"
	"github.com/rs/zerolog"
)

type CreatePollCommandHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewCreatePollCommandHandler(logger zerolog.Logger, repository Repository) CreatePollCommandHandler {
	return CreatePollCommandHandler{logger, repository}
}

func (h CreatePollCommandHandler) Handle(request interface{}) (interface{}, error) {
	req, ok := request.(*pollpb.CreatePollRequest)
	if !ok {
		return nil, nil
	}

	entity := entities.Poll{
		CreatorId:          req.CreatorId,
		Title:              req.Title,
		AllowedUniqueVotes: req.AllowedUniqueVotes,
	}

	for _, option := range req.Options {
		optionEntity := entities.Option{
			Content: option.Content,
		}

		entity.Options = append(entity.Options, optionEntity)
	}

	poll, err := h.repository.CreatePoll(entity)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Str("title", poll.Title).
		Str("creatorId", poll.CreatorId).
		Int32("allowedUniqueVotes", poll.AllowedUniqueVotes).
		Msg("poll created")

	pbPoll := mapToPbPoll(poll)

	return &pollpb.CreatePollReply{Poll: pbPoll}, nil
}
