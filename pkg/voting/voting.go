package voting

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
	VoterId  string
	ServerId string
	Options  []VoteOption
}

type VoteOption struct {
	Rank  int32
	Index int32
}

type PollService interface {
	Create(poll Poll) (*Poll, error)
	Poll(id string) (*Poll, error)
	PollCreator(id string) (string, error)
	End(id string) (*Poll, error)
	HasEnded(id string) (bool, error)
}

type SessionService interface {
	CurrentPoll(serverId string) (string, error)
	SetCurrentPoll(serverId, pollId string) error
}

type BallotService interface {
	Create(poll Poll) (*Ballot, error)
	Ballot(serverId, voterId string) (*Ballot, error)
	Submit(vote Vote) error
	Count(pollId string) error
}
