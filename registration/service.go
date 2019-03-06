package registration

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/rs/zerolog"
)

type Service interface {
	RegisterVoter(*registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error)
}

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
}

func NewDefaultService(logger zerolog.Logger, repository Repository) DefaultService {
	return DefaultService{logger, repository}
}

func (h DefaultService) RegisterVoter(req *registrationpb.RegisterVoterRequest) (*registrationpb.RegisterVoterReply, error) {
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
