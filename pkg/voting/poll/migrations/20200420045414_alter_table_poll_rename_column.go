package migrations

import "database/sql"

type AlterTablePollRenameColumn20200420045414 struct{}

func (m AlterTablePollRenameColumn20200420045414) Version() string {
	return "20200420045414_AlterTablePollRenameColumn"
}

func (m AlterTablePollRenameColumn20200420045414) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE poll RENAME COLUMN hasEnded TO manuallyEnded`)
	return err
}

func (m AlterTablePollRenameColumn20200420045414) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE poll RENAME COLUMN manuallyEnded TO hasEnded`)
	return err
}
