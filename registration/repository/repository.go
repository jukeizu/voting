package repository

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/registration/repository/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "registration"
)

type Repository interface {
	Migrate() error
	RegisterVoter(externalId, username string, canVote bool) (*registrationpb.Voter, error)
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

	err = g.RegisterMigrations(migrations.CreateTableVoter20190221024754{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) RegisterVoter(externalId string, username string, canVote bool) (*registrationpb.Voter, error) {
	q := `INSERT into voter (externalId, username, canVote)
	VALUES ($1, $2, $3)
	ON CONFLICT (externalId) 
	DO UPDATE SET username = excluded.username, updated = NOW()
	RETURNING id, externalId, username, canvote`

	voter := registrationpb.Voter{}

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
