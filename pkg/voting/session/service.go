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
