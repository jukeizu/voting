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
	Poll(shortId string, serverId string) (*voting.Poll, error)
	HasEnded(shortId string, serverId string) (bool, error)
	PollCreator(shortId string, serverId string) (string, error)
	Options(pollId string) ([]voting.Option, error)
	EndPoll(pollShortId, serverId string) (*voting.Poll, error)
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
	q := `INSERT INTO poll (
		shortId,
		serverId,
		creatorId, 
		title, 
		allowedUniqueVotes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	poll := req

	err := r.Db.QueryRow(q,
		req.ShortId,
		req.ServerId,
		req.Title,
		req.CreatorId,
		req.AllowedUniqueVotes,
	).Scan(&poll.Id)
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

func (r *repository) Poll(shortId string, serverId string) (*voting.Poll, error) {
	q := `SELECT id, 
		shortId,
		serverId,
		creatorId, 
		title, 
		allowedUniqueVotes, 
		hasEnded
	FROM poll WHERE shortid = $1 AND serverid = $2`

	poll := voting.Poll{}

	err := r.Db.QueryRow(q, shortId, serverId).Scan(
		&poll.Id,
		&poll.ShortId,
		&poll.ServerId,
		&poll.CreatorId,
		&poll.Title,
		&poll.AllowedUniqueVotes,
		&poll.HasEnded,
	)
	if err != nil {
		return nil, err
	}

	options, err := r.Options(poll.Id)
	if err != nil {
		return nil, err
	}

	poll.Options = options

	return &poll, nil
}

func (r *repository) HasEnded(shortId string, serverId string) (bool, error) {
	q := `SELECT hasEnded FROM poll WHERE shortid = $1 AND serverid = $2`

	poll := voting.Poll{}

	err := r.Db.QueryRow(q, shortId, serverId).Scan(
		&poll.HasEnded,
	)
	if err != nil {
		return false, err
	}

	return poll.HasEnded, nil
}

func (r *repository) PollCreator(shortId string, serverId string) (string, error) {
	q := `SELECT creatorId FROM poll WHERE shortid = $1 AND serverid = $2`

	creator := ""

	err := r.Db.QueryRow(q, shortId, serverId).Scan(&creator)
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

func (r *repository) EndPoll(pollShortId string, serverId string) (*voting.Poll, error) {
	q := `UPDATE poll SET hasEnded = true WHERE shortId = $1 AND serverId = $2`

	_, err := r.Db.Exec(q, pollShortId, serverId)
	if err != nil {
		return nil, err
	}

	return r.Poll(pollShortId, serverId)
}

func (r *repository) createOption(pollId string, option voting.Option) (*voting.Option, error) {
	q := `INSERT INTO option (pollid, content, url) VALUES ($1, $2, $3) RETURNING id`

	o := option

	err := r.Db.QueryRow(q,
		pollId,
		option.Content,
		option.Url,
	).Scan(
		&o.Id,
	)

	return &o, err
}
