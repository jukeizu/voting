package registration

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/api/protobuf-spec/registrationpb"
	"github.com/jukeizu/voting/registration/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "registration"
)

type Repository interface {
	Migrate() error
	SaveUser(externalId, username string) (*registrationpb.User, error)
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

	err = g.RegisterMigrations(migrations.CreateTableUser20190221024754{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) SaveUser(externalId string, username string) (*registrationpb.User, error) {
	q := `INSERT into user (externalId, username, canVote)
	VALUES ($1, $2)
	ON CONFLICT (externalId) 
	DO UPDATE SET username = excluded.username, updated = NOW()
	RETURNING id, externalId, username, canvote`

	user := registrationpb.User{}

	err := r.Db.QueryRow(q,
		externalId,
		username,
	).Scan(
		&user.Id,
		&user.ExternalId,
		&user.Username,
		&user.CanVote,
	)

	return &user, err
}
