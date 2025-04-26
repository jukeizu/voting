package migrations

import "database/sql"

type CreateTableCandidate20250105064414 struct{}

func (m CreateTableCandidate20250105064414) Version() string {
	return "20250105064414_CreateTableOption"
}

func (m CreateTableCandidate20250105064414) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS candidate (
			id SERIAL PRIMARY KEY,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			poll_id integer NOT NULL,
			name TEXT NOT NULL DEFAULT '',
			url TEXT,
			FOREIGN KEY (poll_id) REFERENCES poll (id) ON DELETE CASCADE
		)`)
	return err
}

func (m CreateTableCandidate20250105064414) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE candidate`)
	return err
}
