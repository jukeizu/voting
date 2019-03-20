package migrations

import "database/sql"

type CreateTableBallotOption20190320020149 struct{}

func (m CreateTableBallotOption20190320020149) Version() string {
	return "20190320020149_CreateTableBallotOption"
}

func (m CreateTableBallotOption20190320020149) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS ballot_option (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			ballotId UUID NOT NULL REFERENCES ballot,
			index INT NOT NULL,
			optionId UUID NOT NULL,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableBallotOption20190320020149) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE ballot`)
	return err
}
