package migrations

import "database/sql"

type CreateTableBallot20250106041949 struct{}

func (m CreateTableBallot20250106041949) Version() string {
	return "20250106041949_CreateTableBallotOption"
}

func (m CreateTableBallot20250106041949) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS ballot (
			id SERIAL PRIMARY KEY,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			voter_id integer NOT NULL,
			candidate_id integer NOT NULL,
			rank SMALLINT NOT NULL DEFAULT 1,
			void BOOL NOT NULL DEFAULT false,
			FOREIGN KEY (voter_id) REFERENCES voter (id),
			FOREIGN KEY (candidate_id) REFERENCES candidate (id)
		)`)
	return err
}

func (m CreateTableBallot20250106041949) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE ballot`)
	return err
}
