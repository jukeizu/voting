package voting

import (
	"github.com/rs/zerolog"
)

type Service interface {
	CreatePoll(poll Poll) (*Poll, error)
	Poll(shortId string, serverId string) (*Poll, error)
	EndPoll(shortId string, serverId string, userId string) (*Poll, error)
	Vote(vote Vote) error
	Count(pollId string) error
	CurrentPoll(serverId string) (string, error)
	SetCurrentPoll(serverId, pollId string) error
}

type DefaultService struct {
	logger         zerolog.Logger
	pollService    PollService
	sessionService SessionService
}

func NewDefaultService(
	logger zerolog.Logger,
	pollService PollService,
	sessionService SessionService,
) DefaultService {
	return DefaultService{
		logger,
		pollService,
		sessionService,
	}
}

func (s DefaultService) CreatePoll(poll Poll) (*Poll, error) {
	return s.pollService.Create(poll)
}

func (s DefaultService) Poll(shortId string, serverId string) (*Poll, error) {
	return s.pollService.Poll(shortId, serverId)
}

func (s DefaultService) EndPoll(shortId string, serverId string, userId string) (*Poll, error) {
	return s.pollService.End(shortId, serverId)
}

func (s DefaultService) Vote(vote Vote) error {
	return nil
}

func (s DefaultService) Count(pollId string) error {
	return nil
}

func (s DefaultService) CurrentPoll(serverId string) (string, error) {
	return s.sessionService.CurrentPoll(serverId)
}

func (s DefaultService) SetCurrentPoll(serverId string, pollId string) error {
	return s.sessionService.SetCurrentPoll(serverId, pollId)
}
