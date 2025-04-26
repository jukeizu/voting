package migrations

import "database/sql"

type CreateTableVoter20250105064406 struct{}

func (m CreateTableVoter20250105064406) Version() string {
	return "20250105064406_CreateTableVoter"
}

func (m CreateTableVoter20250105064406) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS voter (
			id SERIAL PRIMARY KEY,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			organization_id integer NOT NULL,
			external_id TEXT NOT NULL,
			name TEXT NOT NULL DEFAULT '',
			can_vote BOOL NOT NULL DEFAULT true,
			UNIQUE(organization_id, external_id),
			FOREIGN KEY (organization_id) REFERENCES organization (id)
		)`)
	return err
}

func (m CreateTableVoter20250105064406) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE voter`)
	return err
}
