package voting

import (
	"fmt"

	"github.com/rs/zerolog"
)

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

func (s ValidationService) Vote(voteRequest VoteRequest) (VoteReply, error) {
	valid, message, err := s.validatePollIsActive(voteRequest.ShortId, voteRequest.ServerId)
	if err != nil {
		return VoteReply{}, err
	}
	if !valid {
		return VoteReply{Message: message}, nil
	}

	poll, err := s.pollService.Poll(voteRequest.ShortId, voteRequest.ServerId)
	if err != nil {
		return VoteReply{}, err
	}

	valid, message, err = s.validateBallotOptions(poll, voteRequest.Options)
	if err != nil {
		return VoteReply{}, err
	}
	if !valid {
		return VoteReply{Message: message}, nil
	}

	return s.service.Vote(voteRequest)
}

func (s ValidationService) Count(pollId string) error {
	return nil
}

func (s ValidationService) CurrentPoll(serverId string) (string, error) {
	return s.service.CurrentPoll(serverId)
}

func (s ValidationService) SetCurrentPoll(serverId string, pollId string) error {
	return s.service.SetCurrentPoll(serverId, pollId)
}

func (s ValidationService) validatePollIsActive(shortId string, serverId string) (bool, string, error) {
	hasEnded, err := s.pollService.HasEnded(shortId, serverId)
	if err != nil {
		return false, "", err
	}

	if hasEnded {
		return false, ErrPollHasEnded.Error(), nil
	}

	return true, "", nil
}

func (s ValidationService) validatePollHasEnded(shortId string, serverId string) (bool, string, error) {
	hasEnded, err := s.pollService.HasEnded(shortId, serverId)
	if err != nil {
		return false, "", err
	}

	if !hasEnded {
		return false, ErrPollHasNotEnded.Error(), nil
	}

	return true, "", nil
}

func (s ValidationService) validateBallotOptions(poll Poll, ballotOptions []BallotOption) (bool, string, error) {
	optionCount := len(ballotOptions)
	if optionCount < 1 {
		return false, "At least one option must be specified.", nil
	}

	if optionCount > int(poll.AllowedUniqueVotes) {
		return false, fmt.Sprintf("Too many votes. Maximum for this poll is %d.", poll.AllowedUniqueVotes), nil
	}

	optionIds := []string{}
	for _, option := range ballotOptions {
		optionIds = append(optionIds, option.OptionId)
	}

	availableOptions, err := s.pollService.UniqueOptions(poll.Id, optionIds)
	if err != nil {
		return false, "", err
	}

	if len(availableOptions) != optionCount {
		return false, "Your vote contains invalid or duplicate options.", nil
	}

	return true, "", nil
}
