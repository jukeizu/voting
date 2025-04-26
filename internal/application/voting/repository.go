package voting

import (
	"database/sql"

	"github.com/jukeizu/voting/internal/application"
	"github.com/jukeizu/voting/internal/database"
	"github.com/lib/pq"
)

type Repository struct {
	db application.Database
}

func NewRepository(db application.Database) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) VoidPreviousBallot(voterID string, pollID string) error {
	q := `UPDATE ballot
		SET void = true, updated = now()
		FROM candidate
		WHERE voter_id = $1
	 	AND candidate.id = candidate_id
		AND candidate.poll_id = $2`

	_, err := r.db.Exec(q, voterID, pollID)

	return err
}

func (r Repository) SaveRankedChoices(voterID string, options []application.RankedChoice) error {
	return database.ExecMultiple(
		r.db.SQL(),
		pq.CopyIn("ballot", "voter_id", "candidate_id", "rank"),
		options,
		func(stmt *sql.Stmt, b application.RankedChoice) (sql.Result, error) {
			return stmt.Exec(voterID, b.Choice.CandidateID, b.Rank)
		},
	)
}

func (r Repository) VoterSession(voterID string) (application.Session, error) {
	q := `SELECT id,
			voter_id,
			poll_id,
			salt
		FROM session
		WHERE voter_id = $1
		ORDER BY last_viewed desc
		LIMIT 1`

	session := application.Session{}
	err := r.db.
		QueryRow(q, voterID).
		Scan(
			&session.ID,
			&session.VoterID,
			&session.PollID,
			&session.Salt,
		)

	if err != nil {
		return application.Session{}, err
	}

	return session, nil
}

func (r Repository) SaveVoterSession(voterID string, pollID string) (application.Session, error) {
	q := `INSERT INTO session (
		voter_id,
		poll_id,
		last_viewed)
	VALUES ($1, $2, now())
	ON CONFLICT (voter_id, poll_id)
	DO UPDATE SET last_viewed = now() 
	RETURNING id, salt`

	session := application.Session{
		PollID:  pollID,
		VoterID: voterID,
	}
	err := r.db.
		QueryRow(q, voterID, pollID).
		Scan(&session.ID, &session.Salt)

	if err != nil {
		return application.Session{}, err
	}

	return session, nil
}

func (r Repository) BallotChoices(session application.Session) ([]application.Choice, error) {
	choices := []application.Choice{}

	q := `SELECT row_number() OVER (ORDER BY md5(concat(candidate.id::text, $2::text, $3::text)) COLLATE "C" DESC) AS number, id, name, url
		FROM candidate
		WHERE poll_id = $1
		ORDER BY number ASC`

	rows, err := r.db.Query(q, session.PollID, session.VoterID, session.Salt)
	if err != nil {
		return choices, err
	}

	defer rows.Close()
	for rows.Next() {
		choice := application.Choice{}
		err := rows.Scan(
			&choice.Number,
			&choice.CandidateID,
			&choice.Name,
			&choice.URL,
		)
		if err != nil {
			return choices, err
		}

		choices = append(choices, choice)
	}

	return choices, nil
}

func (r Repository) Ballot(session application.Session) (application.Ballot, error) {
	ballot := application.Ballot{
		Voter:  application.Voter{ID: session.VoterID},
		PollId: session.PollID,
	}

	q := `SELECT rank, c.number, c.id, c.name, c.url
		FROM ballot
		INNER JOIN (SELECT row_number() OVER (ORDER BY md5(concat(candidate.id::text, $2::text, $3::text)) COLLATE "C" DESC) AS number, id, name, url FROM candidate WHERE poll_id = $1) c 
			ON candidate_id = c.id
		WHERE voter_id = $2::integer
		AND void = false
		ORDER BY rank`

	rows, err := r.db.Query(q, session.PollID, session.VoterID, session.Salt)
	if err != nil {
		return ballot, err
	}

	defer rows.Close()
	for rows.Next() {
		option := application.RankedChoice{}
		err := rows.Scan(
			&option.Rank,
			&option.Choice.Number,
			&option.Choice.CandidateID,
			&option.Choice.Name,
			&option.Choice.URL,
		)
		if err != nil {
			return ballot, err
		}

		ballot.RankedChoices = append(ballot.RankedChoices, option)
	}

	return ballot, nil
}
