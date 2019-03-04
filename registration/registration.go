package registration

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/rs/zerolog"
)

type RegistrationHandler interface {
	RegisterVoter(*registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error)
}

type registrationHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewRegistrationHandler(logger zerolog.Logger, repository Repository) RegistrationHandler {
	return &registrationHandler{logger, repository}
}

func (h registrationHandler) RegisterVoter(req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
	voter, err := h.repository.RegisterVoter(req.ExternalId, req.Username, true)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("externalId", voter.ExternalId).
		Str("username", voter.Username).
		Msg("registered voter")

	return &registrationpb.RegisterVoterReply{Voter: voter}, nil
}
