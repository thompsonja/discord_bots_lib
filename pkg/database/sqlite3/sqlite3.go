package sqlite3

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabaseWithMigrations(dbfile, migrationsDir string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create sqlite3 driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"sqlite3", driver)
	if err != nil {
		return nil, fmt.Errorf("could not create new migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("failed to run migrate (%v): %v", "file://"+migrationsDir, err)
	}

	return db, nil
}
