package voting

type Poll struct {
	Id                 string
	ShortId            string
	ServerId           string
	CreatorId          string
	Title              string
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

type Vote struct {
	PollId   string
	ServerId string
	Voter    Voter
	Options  []VoteOption
}

type VoteOption struct {
	Rank     int32
	OptionId string
}

type PollService interface {
	Create(poll Poll) (Poll, error)
	Poll(shortId string, serverId string) (Poll, error)
	PollCreator(shortId string, serverId string) (string, error)
	End(shortId string, serverId string) (Poll, error)
	HasEnded(shortId string, serverId string) (bool, error)
}

type SessionService interface {
	CurrentPoll(serverId string) (string, error)
	SetCurrentPoll(serverId, pollId string) error
}
