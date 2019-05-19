package migrations

import "database/sql"

type CreateTableVoterPoll20190519234636 struct{}

func (m CreateTableVoterPoll20190519234636) Version() string {
	return "20190519234636_CreateTableVoterPoll"
}

func (m CreateTableVoterPoll20190519234636) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS voter_poll (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			voterId STRING NOT NULL DEFAULT '',
			serverId STRING UNIQUE NOT NULL DEFAULT '',
			pollId STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableVoterPoll20190519234636) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE voter_poll`)
	return err
}
