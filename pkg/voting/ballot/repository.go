package ballot

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/pkg/voting"
	"github.com/jukeizu/voting/pkg/voting/ballot/migrations"
	"github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "voting_ballot"
)

type Repository interface {
	Migrate() error
	VoidBallotOptions(pollId string, voterId string) error
	CreateBallotOptions(ballot voting.Ballot) error
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

	err = g.RegisterMigrations(migrations.CreateTableBallot20190423035611{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) VoidBallotOptions(pollId string, voterId string) error {
	q := `UPDATE ballot_option SET void = true WHERE pollId = $1 AND voterId = $2`

	_, err := r.Db.Exec(q, pollId, voterId)

	return err
}

func (r *repository) CreateBallotOptions(ballot voting.Ballot) error {
	txn, err := r.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("ballot_option", "pollid", "optionid", "rank", "voterid"))
	if err != nil {
		return err
	}

	for _, ballotOption := range ballot.Options {
		_, err := stmt.Exec(ballot.PollId, ballotOption.OptionId, ballotOption.Rank, ballot.Voter.Id)
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
