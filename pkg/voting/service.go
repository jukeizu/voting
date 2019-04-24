package voting

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/election"
	"github.com/shawntoffel/electioncounter"
)

type Service interface {
	CreatePoll(poll Poll) (Poll, error)
	Poll(shortId string, serverId string) (Poll, error)
	EndPoll(shortId string, serverId string, userId string) (Poll, error)
	Status(shortId string, serverId string) (Status, error)
	Vote(voteRequest VoteRequest) (VoteReply, error)
	Count(countRequest CountRequest) (countResult CountResult, err error)
	CurrentPoll(serverId string) (string, error)
	SetCurrentPoll(serverId, pollId string) error
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
	return s.pollService.Create(poll)
}

func (s DefaultService) Poll(shortId string, serverId string) (Poll, error) {
	return s.pollService.Poll(shortId, serverId)
}

func (s DefaultService) EndPoll(shortId string, serverId string, userId string) (Poll, error) {
	return s.pollService.End(shortId, serverId)
}

func (s DefaultService) Status(shortId string, serverId string) (Status, error) {
	poll, err := s.pollService.Poll(shortId, serverId)
	if err != nil {
		return Status{}, err
	}

	status := Status{
		Poll: poll,
	}

	voterIds, err := s.ballotService.VoterIds(poll.Id)
	if err != nil {
		return Status{}, err
	}

	for _, voterId := range voterIds {
		voter, err := s.voterService.Voter(voterId)
		if err != nil {
			return Status{}, err
		}

		status.Voters = append(status.Voters, voter)
	}

	return status, nil
}

func (s DefaultService) Vote(voteRequest VoteRequest) (VoteReply, error) {
	poll, err := s.pollService.Poll(voteRequest.ShortId, voteRequest.ServerId)
	if err != nil {
		return VoteReply{}, errors.New("couldn't find poll: " + err.Error())
	}

	voteRequest.Voter.CanVote = true

	voter, err := s.voterService.Create(voteRequest.Voter)
	if err != nil {
		return VoteReply{}, errors.New("couldn't find voter: " + err.Error())
	}

	ballot := Ballot{
		PollId:  poll.Id,
		Voter:   voter,
		Options: voteRequest.Options,
	}

	ballotResult, err := s.ballotService.Submit(ballot)
	if err != nil {
		return VoteReply{}, errors.New("couldn't submit ballot: " + err.Error())
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
	poll, err := s.pollService.Poll(countRequest.ShortId, countRequest.ServerId)
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
		Precision:           6,
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

func (s DefaultService) CurrentPoll(serverId string) (string, error) {
	return s.sessionService.CurrentPoll(serverId)
}

func (s DefaultService) SetCurrentPoll(serverId string, pollId string) error {
	return s.sessionService.SetCurrentPoll(serverId, pollId)
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
