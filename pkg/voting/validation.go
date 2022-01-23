package voting

import (
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type ValidationService struct {
	logger       zerolog.Logger
	service      Service
	pollService  PollService
	voterService VoterService
}

func NewValidationService(
	logger zerolog.Logger,
	service Service,
	pollService PollService,
	voterService VoterService,
) Service {
	return &ValidationService{
		logger,
		service,
		pollService,
		voterService,
	}
}

func (s ValidationService) CreatePoll(poll Poll) (Poll, error) {
	if len(poll.Options) < 1 {
		return Poll{}, ErrNoOptions

	}

	if !poll.Expires.IsZero() && poll.Expires.UTC().Before(time.Now().UTC()) {
		return Poll{}, ErrPastPollExpiration
	}

	return s.service.CreatePoll(poll)
}

func (s ValidationService) Poll(shortId string, voterId string, serverId string) (Poll, error) {
	return s.service.Poll(shortId, voterId, serverId)
}

func (s ValidationService) VoterPoll(voterId string, serverId string) (Poll, error) {
	return s.service.VoterPoll(voterId, serverId)
}

func (s ValidationService) EndPoll(shortId string, serverId string, userId string) (Poll, error) {
	poll, err := s.service.Poll(shortId, "", serverId)
	if err != nil {
		return Poll{}, err
	}

	if userId != poll.CreatorId {
		return Poll{}, ErrNotOwner
	}

	return s.service.EndPoll(shortId, serverId, userId)
}

func (s ValidationService) OpenPoll(shortId string, serverId string, userId string, expires time.Time) (OpenPollResult, error) {
	poll, err := s.service.Poll(shortId, "", serverId)
	if err != nil {
		return OpenPollResult{}, err
	}

	if userId != poll.CreatorId {
		return OpenPollResult{}, ErrNotOwner
	}

	return s.service.OpenPoll(shortId, serverId, userId, expires)
}

func (s ValidationService) Status(shortId string, serverId string) (Status, error) {
	return s.service.Status(shortId, serverId)
}

func (s ValidationService) Voters(shortId string, serverId string) ([]Voter, error) {
	return s.service.Voters(shortId, serverId)
}

func (s ValidationService) Vote(voteRequest VoteRequest) (VoteReply, error) {
	if voteRequest.Voter.ExternalId == "" {
		return VoteReply{}, ErrNoVoterExternalId
	}

	voter, err := s.voterService.Create(voteRequest.Voter)
	if err != nil {
		return VoteReply{}, err
	}

	if !voter.CanVote {
		return VoteReply{}, ErrVoterNotPermitted(voter)
	}

	poll, err := s.service.VoterPoll(voteRequest.Voter.ExternalId, voteRequest.ServerId)
	if err != nil {
		return VoteReply{}, err
	}

	if poll.HasEnded() {
		return VoteReply{}, ErrPollHasEnded
	}

	err = s.validateBallotOptions(poll, voteRequest.Options)
	if err != nil {
		return VoteReply{}, err
	}

	return s.service.Vote(voteRequest)
}

func (s ValidationService) Count(countRequest CountRequest) (CountResult, error) {
	err := s.validatePollHasEnded(countRequest.ShortId, countRequest.ServerId)
	if err != nil {
		return CountResult{}, err
	}

	return s.service.Count(countRequest)
}

func (s ValidationService) Export(exportRequest ExportRequest) (ExportResult, error) {
	if !strings.EqualFold(exportRequest.Method, "blt") {
		return ExportResult{}, ErrUnkownExportMethod(exportRequest.Method)
	}

	return s.service.Export(exportRequest)
}

func (s ValidationService) validatePollHasEnded(shortId string, serverId string) error {
	poll, err := s.service.Poll(shortId, "", serverId)
	if err != nil {
		return err
	}

	if !poll.HasEnded() {
		return ErrPollHasNotEnded
	}

	return nil
}

func (s ValidationService) validateBallotOptions(poll Poll, ballotOptions []BallotOption) error {
	optionCount := len(ballotOptions)
	if optionCount < 1 {
		return ErrNoOptions
	}

	if optionCount > int(poll.AllowedUniqueVotes) {
		return ErrTooManyVotes(poll.AllowedUniqueVotes)
	}

	optionIds := []string{}
	for _, option := range ballotOptions {
		optionIds = append(optionIds, option.OptionId)
	}

	availableOptions, err := s.pollService.UniqueOptions(poll.Id, optionIds)
	if err != nil {
		return err
	}

	if len(availableOptions) != optionCount {
		return ErrInvalidOrDuplicateOptions
	}

	return nil
}
