package postgres

//
//import (
//	"context"
//	"log/slog"
//
//	_ "github.com/jackc/pgx/v5/stdlib"
//	"github.com/jmoiron/sqlx"
//	"github.com/lib/pq"
//	"project-managment/project-service/core"
//)
//
//type DB struct {
//	log  *slog.Logger
//	conn *sqlx.DB
//}
//
//func NewDB(log *slog.Logger, address string) (*DB, error) {
//	db, err := sqlx.Connect("pgx", address)
//	if err != nil {
//		log.Error("connection problem", "address", address, "error", err)
//		return nil, err
//	}
//	return &DB{
//		log:  log,
//		conn: db,
//	}, nil
//}
//
//func (db *DB)
//
//
//func (db *DB) Close() error {
//	return db.conn.Close()
//}
