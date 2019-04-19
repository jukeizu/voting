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

func (s DefaultService) Submit(ballot voting.Ballot) (voting.BallotResult, error) {
	return voting.BallotResult{}, nil
}

func (s DefaultService) Count() error {
	return nil
}
