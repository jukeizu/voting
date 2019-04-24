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

func (s DefaultService) VoterIds(pollId string) ([]string, error) {
	return s.repository.VoterIds(pollId)
}

func (s DefaultService) Count(countRequest voting.CountRequest) (voting.CountResult, error) {
	/*
		ballots, err := s.voterBallots(countRequest.PollId)
		if err != nil {
			return voting.CountResult{}, err
		}

		config := election.Config{
			Ballots:             ballots,
			Precision:           8,
			Seed:                1,
			NumSeats:            countRequest.NumToElect,
			WithdrawnCandidates: countRequest.ToExclude,
		}

		electionCounter := electioncounter.NewElectionCounter()

		result, err := electionCounter.Count(countRequest.Method, config)
		if err != nil {
			return voting.CountResult{}, err
		}
	*/
	return voting.CountResult{}, nil
}

func (s DefaultService) VoterBallot(pollId string, voterId string) ([]string, error) {
	return s.repository.VoterBallot(pollId, voterId)
}

/*
func (s DefaultService) voterBallots(pollId string) (election.Ballots, error) {
	ballots := election.Ballots{}

	voterIds, err := s.repository.VoterIds(pollId)
	if err != nil {
		return ballots, err
	}

	for _, voterId := range voterIds {
		ballot := election.NewBallot()

		options, err := s.repository.VoterBallot(pollId, voterId)
		if err != nil {
			return ballots, err
		}

		for _, option := range options {
			ballot.PushBack(option)
		}

		ballots = append(ballots, ballot)
	}

	return ballots, nil
}*/
