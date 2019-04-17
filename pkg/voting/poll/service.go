package poll

import (
	"errors"

	"github.com/jukeizu/voting/pkg/voting"
	"github.com/rs/zerolog"
	"github.com/teris-io/shortid"
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
	shortId, err := shortid.Generate()
	if err != nil {
		return nil, errors.New("could not create a short id: " + err.Error())
	}

	req.ShortId = shortId

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

func (h DefaultService) Poll(shortId string, serverId string) (*voting.Poll, error) {
	return h.repository.Poll(shortId, serverId)
}

func (h DefaultService) PollCreator(shortId string, serverId string) (string, error) {
	return h.repository.PollCreator(shortId, serverId)
}

func (h DefaultService) HasEnded(shortId string, serverId string) (bool, error) {
	return h.repository.HasEnded(shortId, serverId)
}

func (h DefaultService) End(shortId string, serverId string) (*voting.Poll, error) {
	poll, err := h.repository.EndPoll(shortId, serverId)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Msg("poll has ended")

	return poll, nil
}
