package migrations

import "database/sql"

type CreateTableSession20250105064426 struct{}

func (m CreateTableSession20250105064426) Version() string {
	return "20250105064426_CreateTableVoterSession"
}

func (m CreateTableSession20250105064426) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS session (
			id SERIAL PRIMARY KEY,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			last_viewed TIMESTAMPTZ,
			voter_id integer NOT NULL,
			poll_id integer NOT NULL,
			salt UUID NOT NULL DEFAULT gen_random_uuid(),
			UNIQUE (voter_id, poll_id),
			FOREIGN KEY (voter_id) REFERENCES voter (id),
			FOREIGN KEY (poll_id) REFERENCES poll (id)
		)`)
	return err
}

func (m CreateTableSession20250105064426) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE session`)
	return err
}
