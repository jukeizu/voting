package migrations

import (
	"database/sql"
)

type CreateTableVoter20190221024754 struct{}

func (m CreateTableVoter20190221024754) Version() string {
	return "20190221024754_CreateTableVoter"
}

func (m CreateTableVoter20190221024754) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS voter (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			externalId STRING UNIQUE NOT NULL DEFAULT '',
			username STRING NOT NULL DEFAULT '',
			canvote bool NOT NULL DEFAULT false,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableVoter20190221024754) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE voter`)
	return err
}
