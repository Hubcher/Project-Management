package postgres

import (
	"context"
	"log/slog"

	"github.com/Hubcher/project-management/auth-service/internal/core"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func New(log *slog.Logger, address string) (*DB, error) {

	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}
	return &DB{
		log:  log,
		conn: db,
	}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) SaveUser(ctx context.Context, email string, passHash []byte) (uid string, err error) {
	const op = "postgres.saveUser"
	stmt, err := db.conn.Prepare("INSERT INTO users (email, pass_hash) VALUES ($1, $2)")
	if err != nil {
		return "", err
	}
}

func User(ctx context.Context, email string) (core.User, error) {

}

func isAdmin(ctx context.Context, userID string) (bool, error) {

}
