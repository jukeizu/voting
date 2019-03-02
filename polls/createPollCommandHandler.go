package polls

import (
	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/domain/entities"
	"github.com/jukeizu/voting/persistence"
	"github.com/rs/zerolog"
)

type CreatePollCommandHandler struct {
	logger     zerolog.Logger
	repository persistence.Repository
}

func NewCreatePollCommandHandler(logger zerolog.Logger, repository persistence.Repository) CreatePollCommandHandler {
	return CreatePollCommandHandler{logger, repository}
}

func (h CreatePollCommandHandler) Handle(req *pollpb.CreatePollRequest) (*pollpb.CreatePollReply, error) {
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

	pbPoll := mapToPb(poll)

	return &pollpb.CreatePollReply{Poll: pbPoll}, nil
}
