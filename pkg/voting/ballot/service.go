package ballot

import (
	"github.com/jukeizu/voting/pkg/voting"
	"github.com/rs/zerolog"
)

var _ voting.BallotService = &DefaultService{}

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
}

func NewDefaultService(logger zerolog.Logger, repository Repository) DefaultService {
	return DefaultService{logger, repository}
}

func (s DefaultService) Create(poll voting.Poll) (*voting.Ballot, error) {
	return nil, nil
}

func (s DefaultService) Ballot(serverId string, voterId string) (*voting.Ballot, error) {
	return nil, nil
}

func (s DefaultService) Submit(vote voting.Vote) error {
	return nil
}

func (s DefaultService) Count(pollId string) error {
	return nil
}
