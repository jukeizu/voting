package voting

import "time"

type Poll struct {
	Id                 string
	ShortId            string
	ServerId           string
	CreatorId          string
	Title              string
	AllowedUniqueVotes int32
	Expires            time.Time
	ManuallyEnded      bool
	Options            []Option
}

func (p Poll) HasEnded() bool {
	return p.ManuallyEnded || (!p.Expires.IsZero() && p.Expires.Before(time.Now().UTC()))
}

type Option struct {
	Id      string
	PollId  string
	Content string
	Url     string
}

type Voter struct {
	Id         string
	ExternalId string
	Username   string
	CanVote    bool
}

type VoteRequest struct {
	ShortId  string
	ServerId string
	Voter    Voter
	Options  []BallotOption
}

type VoteReply struct {
	Success bool
	Message string
	Options []VoteReplyOption
}

type Ballot struct {
	PollId  string
	Voter   Voter
	Options []BallotOption
}

type BallotOption struct {
	Rank     int32
	OptionId string
}

type VoteReplyOption struct {
	Rank   int32
	Option Option
}

type BallotResult struct {
	Success bool
	Message string
}

type Status struct {
	Poll       Poll
	VoterCount int64
}

type CountRequest struct {
	ShortId    string
	ServerId   string
	NumToElect int
	Method     string
	ToExclude  []string
}

type CountResult struct {
	Success   bool
	Message   string
	Poll      Poll
	Method    string
	Elected   []VoteReplyOption
	Events    []CountEvent
	Summaries []CountEvent
}

type CountEvent struct {
	Description string
}

type PollService interface {
	Create(poll Poll) (Poll, error)
	Poll(shortId string, serverId string) (Poll, error)
	PollCreator(shortId string, serverId string) (string, error)
	End(shortId string, serverId string) (Poll, error)
	UniqueOptions(pollId string, optionIds []string) ([]Option, error)
	Option(id string) (Option, error)
}

type SessionService interface {
	CurrentPoll(serverId string) (string, error)
	SetCurrentPoll(serverId, pollId string) error
	VoterPoll(voterId string, serverId string) (string, error)
	SetVoterPoll(voterId string, serverId string, pollId string) error
}

type VoterService interface {
	Create(voter Voter) (Voter, error)
	Voter(id string) (Voter, error)
	Voters(ids []string) ([]Voter, error)
}

type BallotService interface {
	Submit(Ballot) (BallotResult, error)
	VoterCount(pollId string) (int64, error)
	VoterIds(pollId string) ([]string, error)
	VoterBallot(pollId string, voterId string) ([]string, error)
}
