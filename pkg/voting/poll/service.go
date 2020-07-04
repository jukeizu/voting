package poll

import (
	"database/sql"
	"errors"
	"time"

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

func (h DefaultService) Create(req voting.Poll) (voting.Poll, error) {
	shortId, err := shortid.Generate()
	if err != nil {
		return voting.Poll{}, errors.New("could not create a short id: " + err.Error())
	}

	req.ShortId = shortId

	poll, err := h.repository.CreatePoll(req)
	if err != nil {
		return voting.Poll{}, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Str("shortId", poll.ShortId).
		Str("serverId", poll.ServerId).
		Str("title", poll.Title).
		Str("creatorId", poll.CreatorId).
		Int32("allowedUniqueVotes", poll.AllowedUniqueVotes).
		Msg("poll created")

	return poll, nil
}

func (h DefaultService) Poll(shortId string, serverId string) (voting.Poll, error) {
	poll, err := h.repository.Poll(shortId, serverId)
	if err == sql.ErrNoRows {
		return voting.Poll{}, voting.ErrPollNotFound(shortId)
	}
	if err != nil {
		return voting.Poll{}, err
	}

	return poll, nil
}

func (h DefaultService) PollCreator(shortId string, serverId string) (string, error) {
	return h.repository.PollCreator(shortId, serverId)
}

func (h DefaultService) End(shortId string, serverId string) (voting.Poll, error) {
	poll, err := h.repository.EndPoll(shortId, serverId)
	if err != nil {
		return voting.Poll{}, err
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Msg("poll has ended")

	return poll, nil
}

func (h DefaultService) Open(shortId string, serverId string, expires time.Time) (voting.OpenPollResult, error) {
	previous, err := h.Poll(shortId, serverId)
	if err != nil {
		return voting.OpenPollResult{}, err
	}

	poll, err := h.repository.OpenPoll(shortId, serverId, expires)
	if err != nil {
		return voting.OpenPollResult{}, err
	}

	openPollResult := voting.OpenPollResult{
		Poll:               poll,
		PreviouslyEnded:    previous.HasEnded(),
		PreviousExpiration: previous.Expires,
	}

	h.logger.Info().
		Str("pollId", poll.Id).
		Bool("previouslyEnded", openPollResult.PreviouslyEnded).
		Time("previousExpiration", openPollResult.PreviousExpiration).
		Time("newExpires", poll.Expires).
		Msg("poll has been opened")

	return openPollResult, nil
}

func (h DefaultService) UniqueOptions(pollId string, optionIds []string) ([]voting.Option, error) {
	return h.repository.UniqueOptions(pollId, optionIds)
}

func (h DefaultService) Option(id string) (voting.Option, error) {
	return h.repository.Option(id)
}
