package application

import (
	"database/sql"
	"fmt"
	"time"
)

type Poll struct {
	Id           string
	Organization Organization
	Creator      Voter
	Title        string
	Created      time.Time
	Expires      *time.Time
	Voters       []Voter
	Ended        bool
}

func (p Poll) ExternalID() string {
	return fmt.Sprintf("%x", p.Created.Unix())
}

type Candidates []Candidate

func (candidates Candidates) Names() []string {
	names := make([]string, len(candidates))

	for i, c := range candidates {
		names[i] = c.Name
	}

	return names
}

type Candidate struct {
	Id   int
	Name string
	URL  string
}

type Organization struct {
	Id                 string
	Name               string
	ExternalId         string
	MaxConcurrentPolls *int
}

type Voter struct {
	ID                  string
	ExternalId          string
	Name                string
	CanVote             bool
	Organization        Organization
	IsOrganizationAdmin bool
}

type Identity struct {
	VoterExternalId        string
	VoterName              string
	OrganizationExternalId string
	OrganizationName       string
}

type Session struct {
	ID        string
	PollID    string
	VoterID   string
	OptionMap []string
	Salt      string
}

type Ballot struct {
	PollId        string
	Voter         Voter
	RankedChoices []RankedChoice
}

func (b Ballot) FlatPreferences() []int {
	preferences := []int{}
	for _, ballotOption := range b.RankedChoices {
		preferences = append(preferences, int(ballotOption.Choice.Number))
	}
	return preferences
}

type BallotOption struct {
	Rank   int32
	Number int32
}

type ElectedCandidate struct {
	Rank int32
	Name string
}

type RankedChoice struct {
	Rank   int32
	Choice Choice
}

type Choice struct {
	Number      int32
	CandidateID string
	Name        string
	URL         string
}

type Database interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	SQL() *sql.DB
}
