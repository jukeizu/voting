package session

import (
	"github.com/jukeizu/voting/pkg/voting"
	"github.com/rs/zerolog"
)

var _ voting.SessionService = &DefaultService{}

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
}

func NewDefaultService(logger zerolog.Logger, repository Repository) DefaultService {
	return DefaultService{logger, repository}
}

func (s DefaultService) CurrentPoll(serverId string) (string, error) {
	return s.repository.CurrentPoll(serverId)
}

func (s DefaultService) SetCurrentPoll(serverId string, pollId string) error {
	s.logger.Info().
		Str("serverId", serverId).
		Str("pollId", pollId).
		Msg("setting current poll")

	return s.repository.SetCurrentPoll(serverId, pollId)
}

func (s DefaultService) VoterPoll(voterId string, serverId string) (string, error) {
	return s.repository.VoterPoll(voterId, serverId)
}

func (s DefaultService) SetVoterPoll(voterId string, serverId string, pollId string) error {
	s.logger.Info().
		Str("voterId", voterId).
		Str("serverId", serverId).
		Str("pollId", pollId).
		Msg("setting voter poll")

	return s.repository.SetVoterPoll(voterId, serverId, pollId)
}
