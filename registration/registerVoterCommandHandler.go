package registration

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/rs/zerolog"
)

type RegisterVoterReply interface {
	Handle(*registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error)
}

type RegisterVoterCommandHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func NewRegisterVoterCommandHandler(logger zerolog.Logger, repository Repository) RegisterVoterCommandHandler {
	return RegisterVoterCommandHandler{logger, repository}
}

func (h RegisterVoterCommandHandler) Handle(req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
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
