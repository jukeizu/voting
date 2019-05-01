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
	err := s.repository.VoidBallotOptions(ballot.PollId, ballot.Voter.Id)
	if err != nil {
		return voting.BallotResult{Success: false}, err
	}

	err = s.repository.CreateBallotOptions(ballot)
	if err != nil {
		return voting.BallotResult{Success: false}, err
	}

	return voting.BallotResult{Success: true}, nil
}

func (s DefaultService) VoterCount(pollId string) (int64, error) {
	return s.repository.VoterCount(pollId)
}

func (s DefaultService) VoterIds(pollId string) ([]string, error) {
	return s.repository.VoterIds(pollId)
}

func (s DefaultService) VoterBallot(pollId string, voterId string) ([]string, error) {
	return s.repository.VoterBallot(pollId, voterId)
}
