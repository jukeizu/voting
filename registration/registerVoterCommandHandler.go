package registration

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/persistence"
	"github.com/rs/zerolog"
)

type RegisterVoterCommandHandler struct {
	logger     zerolog.Logger
	repository persistence.Repository
}

func NewRegisterVoterCommandHandler(logger zerolog.Logger, repository persistence.Repository) RegisterVoterCommandHandler {
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

	reply := registrationpb.RegisterVoterReply{
		Voter: &registrationpb.Voter{
			Id:         voter.Id,
			ExternalId: voter.ExternalId,
			Username:   voter.Username,
			CanVote:    voter.CanVote,
		},
	}

	return &reply, nil
}
