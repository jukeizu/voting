package migrations

import "database/sql"

type CreateTableOption20190220043255 struct{}

func (m CreateTableOption20190220043255) Version() string {
	return "20190220043255_CreateTableOption"
}

func (m CreateTableOption20190220043255) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS option (
			id UUID NOT NULL DEFAULT gen_random_uuid(),
			pollid UUID NOT NULL,
			content STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ,
			PRIMARY KEY (pollid, id),
			FOREIGN KEY (pollid) REFERENCES poll (id) ON DELETE CASCADE
		) INTERLEAVE IN PARENT poll (pollid);`)

	return err
}

func (m CreateTableOption20190220043255) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE option`)
	return err
}
