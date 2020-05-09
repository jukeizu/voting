package migrations

import "database/sql"

type AlterTablePollAddColumnExpires20200420050140 struct{}

func (m AlterTablePollAddColumnExpires20200420050140) Version() string {
	return "20200420050140_AlterTablePollAddColumnExpires"
}

func (m AlterTablePollAddColumnExpires20200420050140) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE poll ADD COLUMN expires TIMESTAMPTZ`)
	return err
}

func (m AlterTablePollAddColumnExpires20200420050140) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE poll DROP COLUMN expires`)
	return err
}
