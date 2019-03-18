package voting

import "github.com/rs/zerolog"

type ValidationService struct {
	logger        zerolog.Logger
	service       Service
	pollService   PollService
	ballotService BallotService
}

func NewValidationService(
	logger zerolog.Logger,
	service Service,
	pollService PollService,
	ballotService BallotService,
) ValidationService {
	return ValidationService{
		logger,
		service,
		pollService,
		ballotService,
	}
}

func (s ValidationService) CreatePoll(poll Poll) (*Poll, error) {
	return s.service.CreatePoll(poll)
}

func (s ValidationService) Poll(id string) (*Poll, error) {
	return s.service.Poll(id)
}

func (s ValidationService) CreateBallot(pollId string, voterId string) (*Ballot, error) {
	err := s.validatePollIsActive(pollId)
	if err != nil {
		return nil, err
	}

	return s.service.CreateBallot(pollId, voterId)
}

func (s ValidationService) Vote(vote Vote) error {
	ballot, err := s.ballotService.Ballot(vote.ServerId, vote.VoterId)
	if err != nil {
		return err
	}

	err = s.validatePollIsActive(ballot.PollId)
	if err != nil {
		return err
	}

	return s.service.Vote(vote)
}

func (s ValidationService) Count(pollId string) error {
	err := s.validatePollHasEnded(pollId)
	if err != nil {
		return err
	}

	return s.service.Count(pollId)
}

func (s ValidationService) CurrentPoll(serverId string) (string, error) {
	return s.service.CurrentPoll(serverId)
}

func (s ValidationService) SetCurrentPoll(serverId string, pollId string) error {
	return s.service.SetCurrentPoll(serverId, pollId)
}

func (s ValidationService) validatePollIsActive(pollId string) error {
	hasEnded, err := s.pollService.HasEnded(pollId)
	if err != nil {
		return err
	}

	if hasEnded {
		return ErrPollHasEnded
	}

	return nil
}

func (s ValidationService) validatePollHasEnded(pollId string) error {
	hasEnded, err := s.pollService.HasEnded(pollId)
	if err != nil {
		return err
	}

	if !hasEnded {
		return ErrPollHasNotEnded
	}

	return nil
}
