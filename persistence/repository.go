package persistence

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/domain/entities"
	"github.com/jukeizu/voting/persistence/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "voting"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(url string) (Repository, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return Repository{}, err
	}

	r := Repository{
		Db: db,
	}

	return r, nil
}

func (r Repository) Migrate() error {
	_, err := r.Db.Exec(`CREATE DATABASE IF NOT EXISTS ` + DatabaseName)
	if err != nil {
		return err
	}

	g, err := gossage.New(r.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(
		migrations.CreateTableVoter20190221024754{},
		migrations.CreateTablePoll20190219024013{},
		migrations.CreateTableOption20190220043255{},
	)
	if err != nil {
		return err
	}

	return g.Up()
}

func (r Repository) RegisterVoter(externalId string, username string, canVote bool) (*entities.Voter, error) {
	q := `INSERT into voter (externalId, username, canVote)
	VALUES ($1, $2, $3)
	ON CONFLICT (externalId) 
	DO UPDATE SET username = excluded.username, updated = NOW()
	RETURNING id, externalId, username, canvote`

	voter := entities.Voter{}

	err := r.Db.QueryRow(q,
		externalId,
		username,
		canVote,
	).Scan(
		&voter.Id,
		&voter.ExternalId,
		&voter.Username,
		&voter.CanVote,
	)

	return &voter, err
}

func (r Repository) CreatePoll(req entities.Poll) (*entities.Poll, error) {
	q := `INSERT INTO poll (title, creatorId, allowedUniqueVotes)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			title,
			creatorId,
			allowedUniqueVotes,
			hasEnded`

	poll := entities.Poll{}

	err := r.Db.QueryRow(q, req.Title, req.CreatorId, req.AllowedUniqueVotes).Scan(
		&poll.Id,
		&poll.Title,
		&poll.CreatorId,
		&poll.AllowedUniqueVotes,
		&poll.HasEnded,
	)
	if err != nil {
		return nil, err
	}

	for _, option := range req.Options {
		o, err := r.createOption(poll.Id, option)
		if err != nil {
			return nil, err
		}

		poll.Options = append(poll.Options, o)
	}

	return &poll, err
}

func (r Repository) Poll(id string) (*entities.Poll, error) {
	q := `SELECT id, 
		title, 
		creatorId, 
		allowedUniqueVotes, 
		hasEnded
	FROM poll WHERE id = $1`

	poll := entities.Poll{}

	err := r.Db.QueryRow(q, id).Scan(
		&poll.Id,
		&poll.Title,
		&poll.CreatorId,
		&poll.AllowedUniqueVotes,
		&poll.HasEnded,
	)
	if err != nil {
		return nil, err
	}

	options, err := r.Options(id)
	if err != nil {
		return nil, err
	}

	poll.Options = options

	return &poll, nil
}

func (r Repository) PollCreator(id string) (string, error) {
	q := `SELECT creatorId FROM poll WHERE id = $1`

	creator := ""

	err := r.Db.QueryRow(q, id).Scan(&creator)
	if err != nil {
		return "", err
	}

	return creator, nil
}

func (r Repository) Options(pollId string) ([]entities.Option, error) {
	options := []entities.Option{}

	q := `SELECT id, pollid, content
		FROM option
		WHERE pollid = $1`

	rows, err := r.Db.Query(q, pollId)
	if err != nil {
		return options, err
	}

	defer rows.Close()
	for rows.Next() {
		option := entities.Option{}
		err := rows.Scan(
			&option.Id,
			&option.PollId,
			&option.Content,
		)
		if err != nil {
			return options, err
		}

		options = append(options, option)
	}

	return options, nil
}

func (r Repository) EndPoll(id string) (*entities.Poll, error) {
	q := `UPDATE poll SET hasEnded = true WHERE id = $1`

	_, err := r.Db.Exec(q, id)
	if err != nil {
		return nil, err
	}

	return r.Poll(id)
}

func (r Repository) createOption(pollId string, option entities.Option) (entities.Option, error) {
	q := `INSERT INTO option (pollid, content) VALUES ($1, $2) RETURNING id, pollid, content`

	o := entities.Option{}

	err := r.Db.QueryRow(q,
		pollId,
		option.Content,
	).Scan(
		&o.Id,
		&o.PollId,
		&o.Content,
	)

	return o, err
}
