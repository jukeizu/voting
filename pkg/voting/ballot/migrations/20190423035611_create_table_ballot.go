package migrations

import "database/sql"

type CreateTableBallot20190423035611 struct{}

func (m CreateTableBallot20190423035611) Version() string {
	return "20190423035611_CreateTableBallot"
}

func (m CreateTableBallot20190423035611) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS ballot_option (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			pollId STRING NOT NULL DEFAULT '',
			optionId STRING NOT NULL DEFAULT '',
			rank INT NOT NULL,
			voterId STRING NOT NULL DEFAULT '',
			void BOOL NOT NULL DEFAULT false,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableBallot20190423035611) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE ballot_option`)
	return err
}
