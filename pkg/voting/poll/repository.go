package poll

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/pkg/voting"
	"github.com/jukeizu/voting/pkg/voting/poll/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "poll"
)

type Repository interface {
	Migrate() error
	CreatePoll(voting.Poll) (*voting.Poll, error)
	Poll(id string) (*voting.Poll, error)
	HasEnded(id string) (bool, error)
	PollCreator(id string) (string, error)
	Options(pollId string) ([]voting.Option, error)
	EndPoll(id string) (*voting.Poll, error)
}

type repository struct {
	Db *sql.DB
}

func NewRepository(url string) (Repository, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	r := repository{
		Db: db,
	}

	return &r, nil
}

func (r *repository) Migrate() error {
	_, err := r.Db.Exec(`CREATE DATABASE IF NOT EXISTS ` + DatabaseName)
	if err != nil {
		return err
	}

	g, err := gossage.New(r.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(
		migrations.CreateTablePoll20190219024013{},
		migrations.CreateTableOption20190220043255{},
	)
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) CreatePoll(req voting.Poll) (*voting.Poll, error) {
	q := `INSERT INTO poll (title, creatorId, allowedUniqueVotes)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			title,
			creatorId,
			allowedUniqueVotes,
			hasEnded`

	poll := voting.Poll{}

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

		poll.Options = append(poll.Options, *o)
	}

	return &poll, err
}

func (r *repository) Poll(id string) (*voting.Poll, error) {
	q := `SELECT id, 
		title, 
		creatorId, 
		allowedUniqueVotes, 
		hasEnded
	FROM poll WHERE id = $1`

	poll := voting.Poll{}

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

func (r *repository) HasEnded(id string) (bool, error) {
	q := `SELECT hasEnded FROM poll WHERE id = $1`

	poll := voting.Poll{}

	err := r.Db.QueryRow(q, id).Scan(
		&poll.HasEnded,
	)
	if err != nil {
		return false, err
	}

	return poll.HasEnded, nil
}

func (r *repository) PollCreator(id string) (string, error) {
	q := `SELECT creatorId FROM poll WHERE id = $1`

	creator := ""

	err := r.Db.QueryRow(q, id).Scan(&creator)
	if err != nil {
		return "", err
	}

	return creator, nil
}

func (r *repository) Options(pollId string) ([]voting.Option, error) {
	options := []voting.Option{}

	q := `SELECT id, pollid, content
		FROM option
		WHERE pollid = $1`

	rows, err := r.Db.Query(q, pollId)
	if err != nil {
		return options, err
	}

	defer rows.Close()
	for rows.Next() {
		option := voting.Option{}
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

func (r *repository) EndPoll(pollId string) (*voting.Poll, error) {
	q := `UPDATE poll SET hasEnded = true WHERE id = $1`

	_, err := r.Db.Exec(q, pollId)
	if err != nil {
		return nil, err
	}

	return r.Poll(pollId)
}

func (r *repository) createOption(pollId string, option voting.Option) (*voting.Option, error) {
	q := `INSERT INTO option (pollid, content) VALUES ($1, $2) RETURNING id, pollid, content`

	o := voting.Option{}

	err := r.Db.QueryRow(q,
		pollId,
		option.Content,
	).Scan(
		&o.Id,
		&o.PollId,
		&o.Content,
	)

	return &o, err
}