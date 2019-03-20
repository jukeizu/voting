package ballot

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/pkg/voting/ballot/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "voting_ballot"
)

type Repository interface {
	Migrate() error
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
		migrations.CreateTableBallot20190320015618{},
		migrations.CreateTableBallotOption20190320020149{},
	)
	if err != nil {
		return err
	}

	return g.Up()
}
