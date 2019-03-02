package registration

import (
	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/rs/zerolog"
)

type RegisterVoterCommandHandler struct {
	logger     zerolog.Logger
	repository Repository
}

func (h RegisterVoterCommandHandler) Handle(request interface{}) (interface{}, error) {
	req, ok := request.(*registrationpb.RegisterVoterRequest)
	if !ok {
		return nil, nil
	}

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
