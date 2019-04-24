package voting

import (
	"errors"

	"github.com/rs/zerolog"
)

type Service interface {
	CreatePoll(poll Poll) (Poll, error)
	Poll(shortId string, serverId string) (Poll, error)
	EndPoll(shortId string, serverId string, userId string) (Poll, error)
	Status(shortIdf string, serverId string) (Status, error)
	Vote(voteRequest VoteRequest) (VoteReply, error)
	Count(pollId string) error
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
) DefaultService {
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

	//TODO: get voters & add to status

	status := Status{
		Poll: poll,
	}

	return status, nil
}

func (s DefaultService) Vote(voteRequest VoteRequest) (VoteReply, error) {
	poll, err := s.pollService.Poll(voteRequest.ShortId, voteRequest.ServerId)
	if err != nil {
		return VoteReply{}, errors.New("couldn't find poll: " + err.Error())
	}

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

func (s DefaultService) Count(pollId string) error {
	return nil
}

func (s DefaultService) CurrentPoll(serverId string) (string, error) {
	return s.sessionService.CurrentPoll(serverId)
}

func (s DefaultService) SetCurrentPoll(serverId string, pollId string) error {
	return s.sessionService.SetCurrentPoll(serverId, pollId)
}
