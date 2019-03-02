package migrations

import "database/sql"

type CreateTablePoll20190219024013 struct{}

func (m CreateTablePoll20190219024013) Version() string {
	return "20190219024013_CreateTablePoll"
}

func (m CreateTablePoll20190219024013) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS poll (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			title STRING NOT NULL DEFAULT '',
			creatorId STRING NOT NULL DEFAULT '',
			allowedUniqueVotes INT NOT NULL DEFAULT 0,
			hasEnded bool NOT NULL DEFAULT false,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTablePoll20190219024013) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE poll`)
	return err
}
