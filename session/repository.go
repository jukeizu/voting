package session

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/session/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "voting_session"
)

type Repository interface {
	Migrate() error
	SetCurrentPoll(serverId, pollId string) error
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
		migrations.CreateTableCurrentPoll20190303044144{},
		migrations.CreateTableVoterSession20190303044855{},
		migrations.CreateTableOption20190303045441{},
	)
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) SetCurrentPoll(serverId string, pollId string) error {
	q := `INSERT INTO current_poll (serverId, pollId)
		VALUES ($1, $2)
		ON CONFLICT (serverId)
		DO UPDATE SET pollId = excluded.pollId, updated = now()`

	_, err := r.Db.Exec(q, serverId, pollId)

	return err
}

func (r *repository) CurrentPoll(serverId string) (string, error) {
	q := `SELECT pollId FROM current_poll WHERE serverId = $1`

	pollId := ""

	err := r.Db.QueryRow(q, serverId).Scan(&pollId)

	return pollId, err
}
