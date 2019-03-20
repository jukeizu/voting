package migrations

import "database/sql"

type CreateTableBallot20190320015618 struct{}

func (m CreateTableBallot20190320015618) Version() string {
	return "20190320015618_CreateTableBallot"
}

func (m CreateTableBallot20190320015618) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS ballot (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			pollId STRING NOT NULL DEFAULT '',
			voterId STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableBallot20190320015618) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE ballot`)
	return err
}
