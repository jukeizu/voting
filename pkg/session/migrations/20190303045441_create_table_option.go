package migrations

import "database/sql"

type CreateTableOption20190303045441 struct{}

func (m CreateTableOption20190303045441) Version() string {
	return "20190303045441_CreateTableOption"
}

func (m CreateTableOption20190303045441) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS option (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			sessionId INT NOT NULL DEFAULT 0,
			voterSessionId STRING NOT NULL DEFAULT '',
			pollId STRING NOT NULL DEFAULT '',
			optionId STRING NOT NULL DEFAULT '',
			content STRING NOT NULL DEFAULT '',
			userId STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)

	return err
}

func (m CreateTableOption20190303045441) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE option`)
	return err
}
