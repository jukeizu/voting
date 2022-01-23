package voting

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/election"
	"github.com/shawntoffel/election/export/blt"
	"github.com/shawntoffel/electioncounter"
)

var MaxStatusVoterCount = int64(50)

type Service interface {
	CreatePoll(poll Poll) (Poll, error)
	Poll(shortId string, voterId string, serverId string) (Poll, error)
	VoterPoll(voterId string, serverId string) (Poll, error)
	EndPoll(shortId string, serverId string, userId string) (Poll, error)
	OpenPoll(shortId string, serverId string, userId string, expires time.Time) (OpenPollResult, error)
	Status(shortId string, serverId string) (Status, error)
	Voters(shortId string, serverId string) ([]Voter, error)
	Vote(voteRequest VoteRequest) (VoteReply, error)
	Count(countRequest CountRequest) (countResult CountResult, err error)
	Export(exportRequest ExportRequest) (ExportResult, error)
}

var _ Service = &DefaultService{}

type DefaultService struct {
	logger         zerolog.Logger
	pollService    PollService
	sessionService SessionService
	voterService   VoterService
	ballotService  BallotService
}

func NewDefaultService(
	logger zerolog.Logger,
	pollService PollService,
	sessionService SessionService,
	voterService VoterService,
	ballotService BallotService,
) Service {
	return DefaultService{
		logger,
		pollService,
		sessionService,
		voterService,
		ballotService,
	}
}

func (s DefaultService) CreatePoll(poll Poll) (Poll, error) {
	poll, err := s.pollService.Create(poll)
	if err != nil {
		return Poll{}, err
	}

	err = s.sessionService.SetCurrentPoll(poll.ServerId, poll.ShortId)
	if err != nil {
		return Poll{}, err
	}

	return poll, nil
}

func (s DefaultService) Poll(shortId string, voterId string, serverId string) (Poll, error) {
	poll, err := s.findPoll(shortId, serverId)
	if err != nil {
		return Poll{}, err
	}

	if voterId != "" {
		err = s.sessionService.SetVoterPoll(voterId, serverId, poll.ShortId)
		if err != nil {
			return Poll{}, err
		}
	}

	return poll, nil
}

func (s DefaultService) VoterPoll(voterId string, serverId string) (Poll, error) {
	pollId, err := s.sessionService.VoterPoll(voterId, serverId)
	if err != nil {
		return Poll{}, err
	}

	return s.findPoll(pollId, serverId)
}

func (s DefaultService) EndPoll(shortId string, serverId string, userId string) (Poll, error) {
	shortId, err := s.findPollShortId(shortId, serverId)
	if err != nil {
		return Poll{}, err
	}

	return s.pollService.End(shortId, serverId)
}

func (s DefaultService) OpenPoll(shortId string, serverId string, userId string, expires time.Time) (OpenPollResult, error) {
	shortId, err := s.findPollShortId(shortId, serverId)
	if err != nil {
		return OpenPollResult{}, err
	}

	return s.pollService.Open(shortId, serverId, expires)
}

func (s DefaultService) Status(shortId string, serverId string) (Status, error) {
	poll, err := s.findPoll(shortId, serverId)
	if err != nil {
		return Status{}, err
	}

	voterCount, err := s.ballotService.VoterCount(poll.Id)
	if err != nil {
		return Status{}, err
	}

	status := Status{
		Poll:       poll,
		VoterCount: voterCount,
	}

	return status, nil
}

func (s DefaultService) Voters(shortId string, serverId string) ([]Voter, error) {
	poll, err := s.findPoll(shortId, serverId)
	if err != nil {
		return []Voter{}, err
	}

	voterIds, err := s.ballotService.VoterIds(poll.Id)
	if err != nil {
		return []Voter{}, err
	}

	voters, err := s.voterService.Voters(voterIds)
	if err != nil {
		return []Voter{}, err
	}

	return voters, nil
}

func (s DefaultService) Vote(voteRequest VoteRequest) (VoteReply, error) {
	poll, err := s.VoterPoll(voteRequest.Voter.ExternalId, voteRequest.ServerId)
	if err != nil {
		return VoteReply{}, err
	}

	voter, err := s.voterService.Create(voteRequest.Voter)
	if err != nil {
		return VoteReply{}, err
	}

	ballot := Ballot{
		PollId:  poll.Id,
		Voter:   voter,
		Options: voteRequest.Options,
	}

	ballotResult, err := s.ballotService.Submit(ballot)
	if err != nil {
		return VoteReply{}, err
	}

	voteReply := VoteReply{
		Success: ballotResult.Success,
		Message: ballotResult.Message,
	}

	for _, ballotOption := range ballot.Options {
		for _, option := range poll.Options {
			if option.Id == ballotOption.OptionId {
				voteReplyOption := VoteReplyOption{
					Rank:   ballotOption.Rank,
					Option: option,
				}

				voteReply.Options = append(voteReply.Options, voteReplyOption)
			}
		}
	}

	return voteReply, nil
}

func (s DefaultService) Count(countRequest CountRequest) (countResult CountResult, err error) {
	poll, err := s.findPoll(countRequest.ShortId, countRequest.ServerId)
	if err != nil {
		return CountResult{}, err
	}

	ballots, err := s.electionBallots(poll.Id)
	if err != nil {
		return CountResult{}, err
	}

	candidates := s.electionCandidates(poll)

	config := election.Config{
		Ballots:             ballots,
		Candidates:          candidates,
		Precision:           8,
		Seed:                1,
		NumSeats:            countRequest.NumToElect,
		WithdrawnCandidates: countRequest.ToExclude,
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic from election counter: %s", debug.Stack())
		}
	}()

	electionCounter := electioncounter.NewElectionCounter()

	result, err := electionCounter.Count(countRequest.Method, config)
	if err != nil {
		return CountResult{}, err
	}

	countResult = CountResult{
		Poll:      poll,
		Method:    countRequest.Method,
		Events:    s.toCountEvents(result.Events),
		Summaries: s.toCountEvents(result.Summaries),
	}

	elected, err := s.toVoteReplyOptions(result.Candidates)
	if err != nil {
		return CountResult{}, err
	}

	countResult.Elected = elected
	countResult.Success = true

	return
}

func (s DefaultService) Export(exportRequest ExportRequest) (ExportResult, error) {
	poll, err := s.findPoll(exportRequest.ShortId, exportRequest.ServerId)
	if err != nil {
		return ExportResult{}, err
	}

	ballots, err := s.electionBallots(poll.Id)
	if err != nil {
		return ExportResult{}, err
	}

	candidates := s.electionCandidates(poll)

	config := election.Config{
		Ballots:             ballots,
		Candidates:          candidates,
		Precision:           8,
		Seed:                1,
		NumSeats:            exportRequest.NumToElect,
		WithdrawnCandidates: exportRequest.ToExclude,
	}

	exporter := blt.Blt{}

	result := ExportResult{
		Content: exporter.Export(config),
	}

	return result, nil
}

func (s DefaultService) findPoll(shortId string, serverId string) (Poll, error) {
	shortId, err := s.findPollShortId(shortId, serverId)
	if err != nil {
		return Poll{}, err
	}

	return s.pollService.Poll(shortId, serverId)
}

func (s DefaultService) findPollShortId(shortId string, serverId string) (string, error) {
	if shortId != "" {
		return shortId, nil
	}

	return s.sessionService.CurrentPoll(serverId)
}

func (s DefaultService) electionBallots(pollId string) (election.Ballots, error) {
	ballots := election.Ballots{}

	voterIds, err := s.ballotService.VoterIds(pollId)
	if err != nil {
		return ballots, err
	}

	for _, voterId := range voterIds {
		ballot := election.NewBallot()

		options, err := s.ballotService.VoterBallot(pollId, voterId)
		if err != nil {
			return ballots, err
		}

		for _, option := range options {
			ballot.PushBack(option)
		}

		ballots = append(ballots, ballot)
	}

	return ballots, nil
}

func (s DefaultService) electionCandidates(poll Poll) election.Candidates {
	candidates := election.Candidates{}

	for _, option := range poll.Options {
		candidate := election.Candidate{
			Id:   option.Id,
			Name: option.Content,
		}

		candidates = append(candidates, candidate)
	}

	return candidates
}

func (s DefaultService) toVoteReplyOptions(elected election.Candidates) ([]VoteReplyOption, error) {
	voteReplyOptions := []VoteReplyOption{}

	for _, candidate := range elected {
		voteReplyOption := VoteReplyOption{
			Rank: int32(candidate.Rank),
		}

		option, err := s.pollService.Option(candidate.Id)
		if err != nil {
			return voteReplyOptions, err
		}

		voteReplyOption.Option = option

		voteReplyOptions = append(voteReplyOptions, voteReplyOption)
	}

	return voteReplyOptions, nil
}

func (s DefaultService) toCountEvents(events election.Events) []CountEvent {
	countEvents := []CountEvent{}

	for _, e := range events {
		countEvent := CountEvent{
			Description: e.Description,
		}

		countEvents = append(countEvents, countEvent)
	}

	return countEvents
}
