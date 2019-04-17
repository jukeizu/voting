package voting

import "github.com/rs/zerolog"

type ValidationService struct {
	logger      zerolog.Logger
	service     Service
	pollService PollService
}

func NewValidationService(
	logger zerolog.Logger,
	service Service,
	pollService PollService,
) Service {
	return &ValidationService{
		logger,
		service,
		pollService,
	}
}

func (s ValidationService) CreatePoll(poll Poll) (Poll, error) {
	return s.service.CreatePoll(poll)
}

func (s ValidationService) Poll(shortId string, serverId string) (Poll, error) {
	return s.service.Poll(shortId, serverId)
}

func (s ValidationService) EndPoll(shortId string, serverId string, userId string) (Poll, error) {
	pollCreator, err := s.pollService.PollCreator(shortId, serverId)
	if err != nil {
		return Poll{}, err
	}

	if userId != pollCreator {
		return Poll{}, ErrNotOwner
	}

	return s.service.EndPoll(shortId, serverId, userId)
}

func (s ValidationService) Status(shortId string, serverId string) (Status, error) {
	return s.service.Status(shortId, serverId)
}

func (s ValidationService) Vote(vote Vote) error {
	/*
		ballot, err := s.ballotService.Ballot(vote.ServerId, vote.VoterId)
		if err != nil {
			return err
		}

		err = s.validatePollIsActive(ballot.PollId)
		if err != nil {
			return err
		}

		return s.service.Vote(vote)
	*/

	return nil
}

func (s ValidationService) Count(pollId string) error {
	/*
		err := s.validatePollHasEnded(pollId)
		if err != nil {
			return err
		}

		return s.service.Count(pollId)

	*/

	return nil
}

func (s ValidationService) CurrentPoll(serverId string) (string, error) {
	return s.service.CurrentPoll(serverId)
}

func (s ValidationService) SetCurrentPoll(serverId string, pollId string) error {
	return s.service.SetCurrentPoll(serverId, pollId)
}

func (s ValidationService) validatePollIsActive(shortId string, serverId string) error {
	hasEnded, err := s.pollService.HasEnded(shortId, serverId)
	if err != nil {
		return err
	}

	if hasEnded {
		return ErrPollHasEnded
	}

	return nil
}

func (s ValidationService) validatePollHasEnded(shortId string, serverId string) error {
	hasEnded, err := s.pollService.HasEnded(shortId, serverId)
	if err != nil {
		return err
	}

	if !hasEnded {
		return ErrPollHasNotEnded
	}

	return nil
}
