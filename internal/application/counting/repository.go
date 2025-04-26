package counting

import (
	"github.com/jukeizu/voting/internal/application"
)

type Repository struct {
	db application.Database
}

func NewRepository(db application.Database) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) VoterIDs(pollId string) ([]string, error) {
	voterIds := []string{}

	q := `SELECT DISTINCT voter_id 
		FROM ballot 
		INNER JOIN candidate ON candidate_id = candidate.Id 
		WHERE candidate.poll_id = $1 
		AND void = false`

	rows, err := r.db.Query(q, pollId)
	if err != nil {
		return voterIds, err
	}

	defer rows.Close()
	for rows.Next() {
		voterId := ""
		err := rows.Scan(
			&voterId,
		)
		if err != nil {
			return voterIds, err
		}

		voterIds = append(voterIds, voterId)
	}

	return voterIds, nil
}

func (r Repository) Ballot(pollId string, voterId string) (application.Ballot, error) {
	ballot := application.Ballot{
		Voter:  application.Voter{ID: voterId},
		PollId: pollId,
	}

	q := `SELECT rank, c.number
		FROM ballot
		INNER JOIN (SELECT row_number() OVER (ORDER by id ASC) as number, id FROM candidate WHERE poll_id = $1) c 
			ON candidate_id = c.id
		WHERE voter_id = $2
		AND void = false
		ORDER BY rank`

	rows, err := r.db.Query(q, pollId, voterId)
	if err != nil {
		return ballot, err
	}

	defer rows.Close()
	for rows.Next() {
		option := application.RankedChoice{}
		err := rows.Scan(
			&option.Rank,
			&option.Choice.Number,
		)
		if err != nil {
			return ballot, err
		}

		ballot.RankedChoices = append(ballot.RankedChoices, option)
	}

	return ballot, nil
}

func (r Repository) Candidates(pollId string) (application.Candidates, error) {
	options := application.Candidates{}

	q := `SELECT id, name, url
		FROM candidate
		WHERE poll_id = $1
		ORDER BY id ASC`

	rows, err := r.db.Query(q, pollId)
	if err != nil {
		return options, err
	}

	defer rows.Close()
	for rows.Next() {
		option := application.Candidate{}
		err := rows.Scan(
			&option.Id,
			&option.Name,
			&option.URL,
		)
		if err != nil {
			return options, err
		}

		options = append(options, option)
	}

	return options, nil
}
