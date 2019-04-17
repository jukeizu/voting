package migrations

import "database/sql"

type CreateTableCurrentPoll20190303044144 struct{}

func (m CreateTableCurrentPoll20190303044144) Version() string {
	return "20190303044144_CreateTableCurrentPoll"
}

func (m CreateTableCurrentPoll20190303044144) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS current_poll (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			pollId STRING NOT NULL DEFAULT '',
			serverId STRING UNIQUE NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableCurrentPoll20190303044144) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE current_poll`)
	return err
}
