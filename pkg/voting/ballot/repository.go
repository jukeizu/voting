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
	VoterCount(pollId string) (int64, error)
	VoterIds(pollId string) ([]string, error)
	VoterBallot(pollId string, voterId string) ([]voting.BallotOption, error)
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
	q := `UPDATE ballot_option SET void = true, updated = now() WHERE pollId = $1 AND voterId = $2`

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

func (r *repository) VoterCount(pollId string) (int64, error) {
	q := `SELECT COUNT(DISTINCT voterid) FROM ballot_option WHERE pollid = $1 AND void = false`

	count := int64(0)

	err := r.Db.QueryRow(q, pollId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) VoterIds(pollId string) ([]string, error) {
	q := `SELECT DISTINCT voterId FROM ballot_option WHERE pollid = $1 AND void = false`

	voterIds := []string{}

	rows, err := r.Db.Query(q, pollId)
	if err != nil {
		return voterIds, err
	}

	defer rows.Close()
	for rows.Next() {
		voterId := ""
		err := rows.Scan(
			&voterId,
		)
		if err != nil {
			return voterIds, err
		}

		voterIds = append(voterIds, voterId)
	}

	return voterIds, nil
}

func (r *repository) VoterBallot(pollId string, voterId string) ([]voting.BallotOption, error) {
	q := `SELECT "rank", optionId FROM ballot_option 
		WHERE pollid = $1 AND voterid = $2 AND void = false
		ORDER by rank`

	ballotOptions := []voting.BallotOption{}

	rows, err := r.Db.Query(q, pollId, voterId)
	if err != nil {
		return ballotOptions, err
	}

	defer rows.Close()
	for rows.Next() {
		var rank int32
		var optionId string
		err := rows.Scan(
			&rank,
			&optionId,
		)
		if err != nil {
			return ballotOptions, err
		}

		ballotOptions = append(ballotOptions, voting.BallotOption{
			Rank:     rank,
			OptionId: optionId,
		})
	}

	return ballotOptions, nil

}
