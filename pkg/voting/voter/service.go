package voter

import (
	"github.com/jukeizu/voting/pkg/voting"
	"github.com/rs/zerolog"
)

var _ voting.VoterService = &DefaultService{}

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
}

func NewDefaultService(logger zerolog.Logger, repository Repository) DefaultService {
	return DefaultService{logger, repository}
}

func (s DefaultService) Create(voter voting.Voter) (voting.Voter, error) {
	voter, err := s.repository.Create(voter)
	if err != nil {
		return voting.Voter{}, err
	}

	s.logger.Info().
		Str("voter.id", voter.Id).
		Str("voter.externalId", voter.ExternalId).
		Str("voter.username", voter.Username).
		Bool("voter.canVote", voter.CanVote).
		Msg("added voter")

	return voter, nil
}

func (s DefaultService) Voter(id string) (voting.Voter, error) {
	return s.repository.Voter(id)
}
