package poll

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/pkg/voting"
	"github.com/jukeizu/voting/pkg/voting/poll/migrations"
	"github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "poll"
)

type Repository interface {
	Migrate() error
	CreatePoll(voting.Poll) (voting.Poll, error)
	Poll(shortId string, serverId string) (voting.Poll, error)
	HasEnded(shortId string, serverId string) (bool, error)
	PollCreator(shortId string, serverId string) (string, error)
	Options(pollId string) ([]voting.Option, error)
	EndPoll(pollShortId, serverId string) (voting.Poll, error)
	UniqueOptions(pollId string, optionIds []string) ([]voting.Option, error)
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

func (r *repository) CreatePoll(req voting.Poll) (voting.Poll, error) {
	q := `INSERT INTO poll (
			shortId,
			serverId,
			title, 
			creatorId, 
			allowedUniqueVotes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			id,
			shortId,
			serverId,
			creatorId,
			title,
			allowedUniqueVotes`

	poll := voting.Poll{}

	err := r.Db.QueryRow(q,
		req.ShortId,
		req.ServerId,
		req.Title,
		req.CreatorId,
		req.AllowedUniqueVotes,
	).Scan(
		&poll.Id,
		&poll.ShortId,
		&poll.ServerId,
		&poll.CreatorId,
		&poll.Title,
		&poll.AllowedUniqueVotes,
	)
	if err != nil {
		return voting.Poll{}, err
	}

	err = r.createOptions(poll.Id, req.Options)
	if err != nil {
		return voting.Poll{}, fmt.Errorf("could not create options for poll: %s", err)
	}

	options, err := r.Options(poll.Id)
	if err != nil {
		return voting.Poll{}, err
	}

	poll.Options = options

	return poll, nil
}

func (r *repository) Poll(shortId string, serverId string) (voting.Poll, error) {
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
		return voting.Poll{}, err
	}

	options, err := r.Options(poll.Id)
	if err != nil {
		return voting.Poll{}, err
	}

	poll.Options = options

	return poll, nil
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

func (r *repository) EndPoll(pollShortId string, serverId string) (voting.Poll, error) {
	q := `UPDATE poll SET hasEnded = true WHERE shortId = $1 AND serverId = $2`

	_, err := r.Db.Exec(q, pollShortId, serverId)
	if err != nil {
		return voting.Poll{}, err
	}

	return r.Poll(pollShortId, serverId)
}

func (r *repository) UniqueOptions(pollId string, optionIds []string) ([]voting.Option, error) {
	options := []voting.Option{}

	q := `SELECT id, pollid, content, url
		FROM option
		WHERE pollId = $1 AND id = ANY($2)`

	rows, err := r.Db.Query(q, pollId, pq.Array(optionIds))
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
			&option.Url,
		)
		if err != nil {
			return options, err
		}

		options = append(options, option)
	}

	return options, nil
}

func (r *repository) createOptions(pollId string, options []voting.Option) error {
	txn, err := r.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("option", "pollid", "content", "url"))
	if err != nil {
		return err
	}

	for _, option := range options {
		_, err := stmt.Exec(pollId, option.Content, option.Url)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
