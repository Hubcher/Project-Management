package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Hubcher/project-management/auth-service/internal/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func NewRepository(log *slog.Logger, address string) (*DB, error) {
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

type accountRow struct {
	UserID            string    `db:"user_id"`
	Email             string    `db:"email"`
	PasswordHash      string    `db:"password_hash"`
	Role              string    `db:"role"`
	IsActive          bool      `db:"is_active"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	PasswordChangedAt time.Time `db:"password_changed_at"`
}

func (a accountRow) toCore() core.Account {
	return core.Account{
		UserID:            a.UserID,
		Email:             a.Email,
		PasswordHash:      a.PasswordHash,
		Role:              core.Role(a.Role),
		IsActive:          a.IsActive,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
		PasswordChangedAt: a.PasswordChangedAt,
	}
}

const (
	createAccountQuery = `
        insert into auth_accounts (user_id, email, password_hash, role, is_active)
        values ($1, $2, $3, $4, $5);
    `
	countAccountsQuery = `
        select count(1)
        from auth_accounts;
    `
	getByEmailQuery = `
        select user_id::text, email, password_hash, role, is_active, created_at, updated_at, password_changed_at
        from auth_accounts
        where email = $1;
    `
	getByUserIDQuery = `
        select user_id::text, email, password_hash, role, is_active, created_at, updated_at, password_changed_at
        from auth_accounts
        where user_id = $1;
    `
	deleteByUserIDQuery = `
        delete from auth_accounts
        where user_id = $1;
    `
)

func (db *DB) CreateAccount(ctx context.Context, acc core.Account) error {
	_, err := db.conn.ExecContext(
		ctx,
		createAccountQuery,
		acc.UserID,
		acc.Email,
		acc.PasswordHash,
		string(acc.Role),
		acc.IsActive,
	)
	if err != nil {
		return fmt.Errorf("create account: %w", err)
	}
	return nil
}

func (db *DB) GetByEmail(ctx context.Context, email string) (core.Account, error) {
	var row accountRow
	if err := db.conn.GetContext(ctx, &row, getByEmailQuery, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Account{}, core.ErrAccountNotFound
		}
		return core.Account{}, fmt.Errorf("get by email: %w", err)
	}
	return row.toCore(), nil
}

func (r *DB) GetByUserID(ctx context.Context, userID string) (core.Account, error) {
	var row accountRow
	if err := r.conn.GetContext(ctx, &row, getByUserIDQuery, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.Account{}, core.ErrAccountNotFound
		}
		return core.Account{}, fmt.Errorf("get by user id: %w", err)
	}
	return row.toCore(), nil
}

func (r *DB) CountAccounts(ctx context.Context) (int, error) {
	var count int
	if err := r.conn.GetContext(ctx, &count, countAccountsQuery); err != nil {
		return 0, fmt.Errorf("count accounts: %w", err)
	}
	return count, nil
}

func (r *DB) DeleteByUserID(ctx context.Context, userID string) error {
	res, err := r.conn.ExecContext(ctx, deleteByUserIDQuery, userID)
	if err != nil {
		return fmt.Errorf("delete by user id: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete by user id affected rows: %w", err)
	}
	if affected == 0 {
		return core.ErrAccountNotFound
	}
	return nil
}
