package poll

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/api/protobuf-spec/pollpb"
	"github.com/jukeizu/voting/poll/migrations"
	_ "github.com/lib/pq"
	"github.com/shawntoffel/gossage"
	"github.com/ventu-io/go-shortid"
)

const (
	DatabaseName = "poll"
)

type Repository interface {
	Migrate() error
	CreatePoll(*pollpb.CreatePollRequest) (*pollpb.Poll, error)
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

func (r *repository) CreatePoll(req *pollpb.CreatePollRequest) (*pollpb.Poll, error) {
	q := `INSERT INTO poll (
			title, 
			creatorId, 
			allowedUniqueVotes)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			title,
			creatorId,
			allowedUniqueVotes,
			hasEnded`

	poll := pollpb.Poll{}

	err := r.Db.QueryRow(q,
		req.Title,
		req.CreatorId,
		req.AllowedUniqueVotes,
	).Scan(
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
		o, err := r.createOption(poll.Id, *option)
		if err != nil {
			return nil, err
		}

		poll.Options = append(poll.Options, o)
	}

	return &poll, err
}

func (r *repository) createOption(pollId string, option pollpb.Option) (*pollpb.Option, error) {
	q := `INSERT INTO option (pollid, content) VALUES ($1, $2) RETURNING id, pollid, content`

	o := pollpb.Option{}

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

func generateShortId() (string, error) {
	return shortid.Generate()
}
