package migrations

import "database/sql"

type CreateTableOrganization20250105064358 struct{}

func (m CreateTableOrganization20250105064358) Version() string {
	return "20250105064358_CreateTableOrganization"
}

func (m CreateTableOrganization20250105064358) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS organization (
			id SERIAL PRIMARY KEY,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			name TEXT NOT NULL,
			external_id TEXT NOT NULL,
			max_concurrent_polls SMALLINT NOT NULL DEFAULT 1,
			UNIQUE (name, external_id)
		)`)
	return err
}

func (m CreateTableOrganization20250105064358) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE organization`)
	return err
}
