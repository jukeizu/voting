package polling

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jukeizu/voting/internal/application"
	"github.com/jukeizu/voting/internal/database"
	"github.com/lib/pq"
)

const MaxVotersReturn = 100

type Repository struct {
	db application.Database
}

func NewRepository(db application.Database) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) SavePoll(req CreatePollRequest) (application.Poll, error) {
	q := `INSERT INTO poll (
			organization_id,
			creator_id, 
			title, 
			expiration)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id string

	err := r.db.
		QueryRow(q, req.Creator.Organization.Id, req.Creator.ID, req.Title, req.Expiration).
		Scan(&id)

	if err != nil {
		return application.Poll{}, err
	}

	err = database.ExecMultiple(
		r.db.SQL(),
		pq.CopyIn("candidate", "poll_id", "name", "url"),
		req.Candidates,
		func(stmt *sql.Stmt, c application.Candidate) (sql.Result, error) {
			return stmt.Exec(id, c.Name, c.URL)
		},
	)
	if err != nil {
		return application.Poll{}, fmt.Errorf("could not create candidates for poll: %s", err)
	}

	return r.Poll(id)
}

func (r Repository) Poll(id string) (application.Poll, error) {
	q := `SELECT poll.id,
			poll.created,
			poll.title,
			poll.expiration,
			COALESCE(poll.expiration <= NOW(), false),
			organization.id,
			organization.name,
			organization.external_id,
			voter.id,
			voter.name,
			voter.external_id
		FROM poll
		INNER JOIN organization
			ON poll.organization_id = organization.id
		INNER JOIN voter
			ON poll.creator_id = voter.id
		WHERE poll.id = $1`

	poll := application.Poll{}

	err := r.db.
		QueryRow(q, id).
		Scan(
			&poll.Id,
			&poll.Created,
			&poll.Title,
			&poll.Expires,
			&poll.Ended,
			&poll.Organization.Id,
			&poll.Organization.Name,
			&poll.Organization.ExternalId,
			&poll.Creator.ID,
			&poll.Creator.Name,
			&poll.Creator.ExternalId,
		)
	if err != nil {
		return application.Poll{}, err
	}

	voters, err := r.Voters(poll.Id, MaxVotersReturn)
	if err != nil {
		return application.Poll{}, err
	}
	poll.Voters = voters

	return poll, nil
}

func (r Repository) Voters(pollID string, limit int) ([]application.Voter, error) {
	voters := []application.Voter{}

	q := `SELECT DISTINCT voter_id, voter.name, voter.external_id FROM ballot b
		INNER JOIN candidate c ON c.id = b.candidate_id AND c.poll_id = $1
		INNER JOIN voter ON voter.id = voter_id
		WHERE b.void = false
		ORDER BY voter.name asc
		LIMIT $2`

	rows, err := r.db.Query(q, pollID, limit)
	if err != nil {
		return voters, err
	}

	defer rows.Close()
	for rows.Next() {
		voter := application.Voter{}
		err := rows.Scan(
			&voter.ID,
			&voter.Name,
			&voter.ExternalId,
		)
		if err != nil {
			return voters, err
		}

		voters = append(voters, voter)
	}

	return voters, err
}

func (r Repository) Open(id string, expires *time.Time) error {
	q := `UPDATE poll SET updated=NOW(), expiration = $1 WHERE id = $2`
	_, err := r.db.Exec(q, expires, id)
	return err
}

func (r Repository) End(id string) error {
	q := `UPDATE poll SET updated=NOW(), expiration=NOW() WHERE id = $1`
	_, err := r.db.Exec(q, id)
	return err
}

func (r Repository) LatestPollID(orgID string) (string, error) {
	q := `SELECT id FROM poll 
		WHERE organization_id = $1
		ORDER by created DESC
		LIMIT 1`

	var id string
	err := r.db.
		QueryRow(q, orgID).
		Scan(
			&id,
		)

	return id, err
}

func (r Repository) PollCreatorID(pollID string) (string, error) {
	q := `SELECT creator_id FROM poll WHERE id = $1`

	var creator_id string
	err := r.db.
		QueryRow(q, pollID).
		Scan(&creator_id)

	return creator_id, err
}

func (r Repository) PollTitle(id string) (string, error) {
	q := `SELECT title FROM poll WHERE id = $1`

	var title string
	err := r.db.
		QueryRow(q, id).
		Scan(&title)

	return title, err
}

func (r Repository) PollHasEnded(id string) (bool, error) {
	q := `SELECT COALESCE(poll.expiration <= NOW(), false) FROM poll WHERE id = $1`

	var ended bool
	err := r.db.
		QueryRow(q, id).
		Scan(&ended)

	return ended, err
}

func (r Repository) CanCreatePoll(orgID string) (bool, error) {
	q := `SELECT COUNT(p.id) < o.max_concurrent_polls FROM organization o 
		LEFT JOIN poll p ON p.organization_id = o.id AND (p.expiration > NOW() OR p.expiration isnull) 
		WHERE o.id = $1
		GROUP BY o.id`

	var canCreate bool
	err := r.db.
		QueryRow(q, orgID).
		Scan(&canCreate)

	if err == sql.ErrNoRows {
		return true, nil
	}

	return canCreate, err
}
