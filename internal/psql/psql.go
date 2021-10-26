package psql

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func NewDb(uri string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", uri)
}

func RunMigrations(db *sqlx.DB, service string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	mPath, err := filepath.Abs(fmt.Sprintf("./internal/%s/migrations", service))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s", mPath),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err.Error() != "no change" && err.Error() != "first : file does not exist" {
		return err
	}
	return nil
}
