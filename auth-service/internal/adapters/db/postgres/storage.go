package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Hubcher/project-management/auth-service/internal/core"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func (db *DB) CreateAccount(ctx context.Context, account core.AuthAccount) error {
	const op = "postgresql.CreateAccount"

	_, err := db.conn.ExecContext(
		ctx,
		createAccountQuery,
		account.UserID,
		account.Email,
		account.PasswordHash,
		account.Role,
		account.IsActive,
		account.CreatedAt,
		account.UpdatedAt,
		account.PasswordChangedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%s: %w", op, core.ErrEmailAlreadyExists)
		}
	}
	return nil
}

func (db *DB) GetAccountByEmail(ctx context.Context, email string) (core.AuthAccount, error) {
	const op = "postgresql.GetAccountByEmail"

	var account core.AuthAccount
	if err := db.conn.GetContext(ctx, &account, getAccountByEmailQuery, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.AuthAccount{}, fmt.Errorf("%s: %w", op, core.ErrAccountNotFound)
		}
		return core.AuthAccount{}, fmt.Errorf("%s: %w", op, err)
	}

	return account, nil
}

func (db *DB) GetAccountByUserId(ctx context.Context, userID string) (core.AuthAccount, error) {
	const op = "postgresql.GetAccountByUserId"
	var account core.AuthAccount
	if err := db.conn.GetContext(ctx, &account, getAccountByUserIDQuery, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.AuthAccount{}, fmt.Errorf("%s: %w", op, core.ErrAccountNotFound)
		}
		return core.AuthAccount{}, fmt.Errorf("%s: %w", op, err)
	}

	return account, nil
}

func (db *DB) DeleteAccountByUserId(ctx context.Context, userID string) error {
	const op = "postgresql.DeleteAccountByUserId"
	_, err := db.conn.ExecContext(ctx, deleteAccountByUserIDQuery, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (db *DB) CreateRefreshSession(ctx context.Context, session core.RefreshSession) error {
	const op = "postgresql.CreateRefreshSession"
	_, err := db.conn.ExecContext(
		ctx,
		CreateRefreshSessionQuery,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.UserAgent,
		session.IP,
		session.ExpiresAt,
		session.CreatedAt,
		session.RevokedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%s: %w", op, core.ErrRefreshSessionAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (db *DB) GetRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (core.RefreshSession, error) {
	const op = "postgres.GetRefreshSessionByTokenHash"

	var session core.RefreshSession
	if err := db.conn.GetContext(ctx, &session, GetRefreshSessionByTokenHashQuery, tokenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.RefreshSession{}, fmt.Errorf("%s: %w", op, core.ErrRefreshSessionNotFound)
		}
		return core.RefreshSession{}, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (db *DB) RevokeRefreshSession(ctx context.Context, sessionID string, revokedAt time.Time) error {
	const op = "postgres.RevokeRefreshSession"

	_, err := db.conn.ExecContext(ctx, RevokeRefreshSessionQuery, sessionID, revokedAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (db *DB) RevokeAllRefreshSessionsByUserID(ctx context.Context, userID string, revokedAt time.Time) error {
	const op = "postgres.RevokeAllRefreshSessionsByUserID"

	_, err := db.conn.ExecContext(ctx, RevokeAllRefreshSessionsByUserIDQuery, userID, revokedAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

const (
	createAccountQuery = `
		insert into auth_account (
			user_id,
		    email,
	        password_hash,
		    role,
		    is_active,
		    created_at,
		    updated_at,
		    password_changed_at
		) values ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	getAccountByEmailQuery = `
		select
		    user_id,
		    email,
		    password_hash,
		    role,
		    is_active,
		    created_at,
		    updated_at,
		    password_changed_at
		from auth_account
		where email = $1
		limit 1
	`
	getAccountByUserIDQuery = `
		select
		    user_id,
		    email,
		    password_hash,
		    role,
		    is_active,
		    created_at,
		    updated_at,
		    password_changed_at
		from auth_account
		where user_id = $1
		limit 1
	`
	deleteAccountByUserIDQuery = `
		delete from auth_account where user_id = $1
	`
	CreateRefreshSessionQuery = `
		insert into refresh_sessions (
		    id,
			user_id,
		    token_hash,
		    user_agent,
		    ip,
		    expires_at,
		    created_at,
		    revoked_at
		) values ($1, $2, $3, $4, $5::inet, $6, $7, $8)                          
	`
	GetRefreshSessionByTokenHashQuery = `
		select
			id,
			user_id,
			token_hash,
			user_agent,
			host(ip) as ip,
			expires_at,
			created_at,
			revoked_at
		from refresh_sessions
		where token_hash = $1
		limit 1
	`
	RevokeRefreshSessionQuery = `
		update refresh_sessions
		set revoked_at = $2
		where id = $1 and revoked_at is null
	`
	RevokeAllRefreshSessionsByUserIDQuery = `
		update refresh_sessions
		set revoked_at = $2
		where user_id = $1 and revoked_at is null
	`
)
