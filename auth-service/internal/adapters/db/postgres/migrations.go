package postgres

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (db *DB) Migrate() error {
	db.log.Debug("running migration")
	files, err := iofs.New(migrationFiles, "migrations") // get migrations from
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db.conn.DB, &pgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", files, "pgx", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	err = m.Up()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			db.log.Error("no new migrations to apply")
			return nil
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	db.log.Debug("migration finished")
	return nil
}
