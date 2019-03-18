package poll

import (
	"github.com/jukeizu/voting"
	"github.com/rs/zerolog"
)

var _ voting.PollService = &DefaultService{}

type DefaultService struct {
	logger     zerolog.Logger
	repository Repository
}

func NewDefaultService(logger zerolog.Logger, repository Repository) DefaultService {
	return DefaultService{logger, repository}
}

func (h DefaultService) Create(req voting.Poll) (*voting.Poll, error) {
	poll, err := h.repository.CreatePoll(req)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Str("title", poll.Title).
		Str("creatorId", poll.CreatorId).
		Int32("allowedUniqueVotes", poll.AllowedUniqueVotes).
		Msg("poll created")

	return poll, nil
}

func (h DefaultService) Poll(id string) (*voting.Poll, error) {
	return h.repository.Poll(id)
}

func (h DefaultService) HasEnded(id string) (bool, error) {
	return h.repository.HasEnded(id)
}

func (h DefaultService) End(pollId string) (*voting.Poll, error) {
	poll, err := h.repository.EndPoll(pollId)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Msg("poll has ended")

	return poll, nil
}
