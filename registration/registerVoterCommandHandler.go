package registration

import (
	"github.com/jukeizu/voting/domain/entities"
	"github.com/jukeizu/voting/persistence"
	"github.com/rs/zerolog"
)

type RegisterVoterRequest struct {
	ExternalId string
	Username   string
}

type RegisterVoterCommandHandler struct {
	logger     zerolog.Logger
	repository persistence.Repository
}

func NewRegisterVoterCommandHandler(logger zerolog.Logger, repository persistence.Repository) RegisterVoterCommandHandler {
	return RegisterVoterCommandHandler{logger, repository}
}

func (h RegisterVoterCommandHandler) Handle(req RegisterVoterRequest) (*entities.Voter, error) {
	voter, err := h.repository.RegisterVoter(req.ExternalId, req.Username, true)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("externalId", voter.ExternalId).
		Str("username", voter.Username).
		Msg("registered voter")

	return voter, nil
}
