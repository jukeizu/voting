package registration

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/rs/zerolog"
)

type RegisterVoterCommandHandler interface {
	Handle(*registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error)
}

type registerVoterCommandHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewRegisterVoterCommandHandler(logger zerolog.Logger, repository Repository) RegisterVoterCommandHandler {
	return &registerVoterCommandHandler{logger, repository}
}

func (h registerVoterCommandHandler) Handle(req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
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
