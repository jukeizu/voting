package migrations

import "database/sql"

type CreateTablePoll20250105064411 struct{}

func (m CreateTablePoll20250105064411) Version() string {
	return "20250105064411_CreateTablePoll"
}

func (m CreateTablePoll20250105064411) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS poll (
			id SERIAL PRIMARY KEY NOT NULL,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			organization_id integer NOT NULL,
			creator_id integer NOT NULL,
			title TEXT NOT NULL DEFAULT '',
			expiration TIMESTAMPTZ,
			FOREIGN KEY (organization_id) REFERENCES organization (id),
			FOREIGN KEY (creator_id) REFERENCES voter (id)
		)`)
	return err
}

func (m CreateTablePoll20250105064411) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE poll`)
	return err
}
