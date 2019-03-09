package session

import (
	"github.com/rs/zerolog"
)

type Service interface {
	CurrentPoll(serverId string) (string, error)
	SetCurrentPoll(serverId, pollId string) error
}

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
