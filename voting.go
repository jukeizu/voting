package main

type Poll struct {
	Id                 string
	Title              string
	CreatorId          string
	AllowedUniqueVotes int32
	HasEnded           bool
	Options            []Option
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

type Ballot struct {
	Id      string
	PollId  string
	VoterId string
	Options []BallotOption
}

type BallotOption struct {
	Id     string
	Index  int32
	Option Option
}

type Vote struct {
	BallotId string
	VoterId  string
	Options  []RankOption
}

type RankOption struct {
	Rank  int32
	Index int32
}

type PollService interface {
	Create(poll Poll) (*Poll, error)
	Poll(id string) (*Poll, error)
	End(id string) (*Poll, error)
}

type BallotService interface {
	Create(poll Poll) (*Ballot, error)
	Submit(vote Vote) error
	Count(pollId string) error
}
