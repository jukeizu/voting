package migrations

import "database/sql"

type CreateTableVoter20190418005651 struct{}

func (m CreateTableVoter20190418005651) Version() string {
	return "20190418005651_CreateTableVoter"
}

func (m CreateTableVoter20190418005651) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS voter (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			externalId STRING NOT NULL DEFAULT '',
			username STRING NOT NULL DEFAULT '',
			canvote BOOL NOT NULL DEFAULT false,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			UNIQUE(externalId)
		)`)

	return err
}

func (m CreateTableVoter20190418005651) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE voter`)
	return err
}
