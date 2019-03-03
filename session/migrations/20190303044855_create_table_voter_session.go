package migrations

import "database/sql"

type CreateTableVoterSession20190303044855 struct{}

func (m CreateTableVoterSession20190303044855) Version() string {
	return "20190303044855_CreateTableVoterSession"
}

func (m CreateTableVoterSession20190303044855) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS voter_session (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			serverId STRING NOT NULL DEFAULT '',
			userId STRING NOT NULL DEFAULT '',
			pollId STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableVoterSession20190303044855) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE voter_session`)
	return err
}
