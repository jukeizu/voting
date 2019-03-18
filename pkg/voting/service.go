package voting

import (
	"github.com/rs/zerolog"
)

type Service interface {
	CreatePoll(poll Poll) (*Poll, error)
	Poll(id string) (*Poll, error)
	EndPoll(id string, userId string) (*Poll, error)
	CreateBallot(pollId, voterId string) (*Ballot, error)
	Vote(vote Vote) error
	Count(pollId string) error
	CurrentPoll(serverId string) (string, error)
	SetCurrentPoll(serverId, pollId string) error
}

type DefaultService struct {
	logger         zerolog.Logger
	pollService    PollService
	ballotService  BallotService
	sessionService SessionService
}

func NewDefaultService(
	logger zerolog.Logger,
	pollService PollService,
	ballotService BallotService,
	sessionService SessionService,
) DefaultService {
	return DefaultService{
		logger,
		pollService,
		ballotService,
		sessionService,
	}
}

func (s DefaultService) CreatePoll(poll Poll) (*Poll, error) {
	return s.pollService.Create(poll)
}

func (s DefaultService) Poll(id string) (*Poll, error) {
	return s.pollService.Poll(id)
}

func (s DefaultService) EndPoll(id string, userId string) (*Poll, error) {
	return s.pollService.End(id)
}

func (s DefaultService) CreateBallot(pollId string, voterId string) (*Ballot, error) {
	poll, err := s.pollService.Poll(pollId)
	if err != nil {
		return nil, err
	}

	return s.ballotService.Create(*poll)
}

func (s DefaultService) Vote(vote Vote) error {
	return s.ballotService.Submit(vote)
}

func (s DefaultService) Count(pollId string) error {
	return s.ballotService.Count(pollId)
}

func (s DefaultService) CurrentPoll(serverId string) (string, error) {
	return s.sessionService.CurrentPoll(serverId)
}

func (s DefaultService) SetCurrentPoll(serverId string, pollId string) error {
	return s.sessionService.SetCurrentPoll(serverId, pollId)
}
