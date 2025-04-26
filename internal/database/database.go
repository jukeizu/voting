package database

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/voting/internal/application"
	"github.com/jukeizu/voting/internal/database/migrations"
	"github.com/lib/pq"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName     = "voting"
	ErrAlreadyExists = "42P04"
	Driver           = "postgres"
)

type Database struct {
	*sql.DB
}

var _ application.Database = &Database{}

func New(url string, migrate bool) (*Database, error) {
	r := &Database{}

	if migrate {
		err := r.createDatabase(url)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open(Driver, fmt.Sprintf(url, DatabaseName))
	if err != nil {
		return nil, err
	}

	r.DB = db

	if migrate {
		err = r.Migrate()
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Database) SQL() *sql.DB {
	return r.DB
}

func (r *Database) Migrate() error {
	g, err := gossage.New(r.DB)
	if err != nil {
		return err
	}
	err = g.RegisterMigrations(
		migrations.CreateTableOrganization20250105064358{},
		migrations.CreateTableVoter20250105064406{},
		migrations.CreateTablePoll20250105064411{},
		migrations.CreateTableCandidate20250105064414{},
		migrations.CreateTableSession20250105064426{},
		migrations.CreateTableBallot20250106041949{},
	)
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *Database) createDatabase(url string) error {
	db, err := sql.Open(Driver, fmt.Sprintf(url, "postgres"))
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE DATABASE ` + DatabaseName)
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if ok && pqErr.Code == ErrAlreadyExists {
		return nil
	}
	return err
}

func ExecMultiple[T any](db *sql.DB, query string, data []T, exec func(stmt *sql.Stmt, d T) (sql.Result, error)) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(query)
	if err != nil {
		return err
	}

	for _, d := range data {
		_, err := exec(stmt, d)
		if err != nil {
			return err
		}
	}

	// Flush buffered data.
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return txn.Commit()
}
