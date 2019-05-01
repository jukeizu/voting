package voter

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/pkg/voting"
	"github.com/jukeizu/voting/pkg/voting/voter/migrations"
	"github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "voting_voter"
)

type Repository interface {
	Migrate() error
	Create(voter voting.Voter) (voting.Voter, error)
	Voter(id string) (voting.Voter, error)
	Voters(ids []string) ([]voting.Voter, error)
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

	err = g.RegisterMigrations(migrations.CreateTableVoter20190418005651{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) Create(voter voting.Voter) (voting.Voter, error) {
	q := `INSERT INTO voter (externalId, username, canVote)
		VALUES ($1, $2, true)
		ON CONFLICT (externalId)
		DO UPDATE SET username = excluded.username, updated = now()
		RETURNING id, canvote`

	err := r.Db.QueryRow(q, voter.ExternalId, voter.Username).Scan(&voter.Id, &voter.CanVote)

	return voter, err
}

func (r *repository) Voter(id string) (voting.Voter, error) {
	q := `SELECT id, externalId, username, canvote FROM voter WHERE id = $1`

	voter := voting.Voter{}

	err := r.Db.QueryRow(q, id).Scan(
		&voter.Id,
		&voter.ExternalId,
		&voter.Username,
		&voter.CanVote,
	)

	return voter, err
}

func (r *repository) Voters(ids []string) ([]voting.Voter, error) {
	q := `SELECT id, externalId, username, canvote FROM voter WHERE id = ANY($1)`

	voters := []voting.Voter{}

	rows, err := r.Db.Query(q, pq.Array(ids))
	if err != nil {
		return voters, err
	}

	defer rows.Close()
	for rows.Next() {
		voter := voting.Voter{}
		err := rows.Scan(
			&voter.Id,
			&voter.ExternalId,
			&voter.Username,
			&voter.CanVote,
		)
		if err != nil {
			return voters, err
		}

		voters = append(voters, voter)
	}

	return voters, nil
}
