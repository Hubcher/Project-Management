package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/Hubcher/project-management/user-service/internal/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

type userRow struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

func (r userRow) toCore() *core.User {
	return &core.User{
		ID:        r.ID,
		Name:      r.Name,
		Password:  r.Password,
		Email:     r.Email,
		Role:      r.Role,
		CreatedAt: r.CreatedAt,
	}
}

func New(log *slog.Logger, address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}
	return &DB{log: log, conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

const (
	createQuery = `
		insert into users (id, name, email, password, role)
		values ($1, $2, $3, $4, $5)
		returning id::text, name, email, password, role, created_at;
	`

	getUserQuery = `
		select id::text, name, email, password, role, created_at
		from users
		where id = $1;
	`

	getAllUsersQuery = `
		select id::text, name, email, password, role, created_at
		from users
		order by created_at desc;
	`

	getUsersByRoleQuery = `
		select id::text, name, email, password, role, created_at
		from users
		where role = $1
		order by created_at desc;
	`

	updateUserQuery = `
		update users
		set name = $2, email = $3, password = $4, role = $5
		where id = $1
		returning id::text, name, email, password, role, created_at;
	`

	deleteUserQuery = `
		delete from users
		where id = $1;
	`
)

func (db *DB) CreateUser(ctx context.Context, input core.CreateUserInput) (*core.User, error) {
	var row userRow
	if err := db.conn.GetContext(
		ctx,
		&row,
		createQuery,
		input.ID,
		input.Name,
		input.Email,
		input.Password,
		input.Role,
	); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return row.toCore(), nil
}

func (db *DB) GetUser(ctx context.Context, id string) (*core.User, error) {
	var row userRow
	if err := db.conn.GetContext(ctx, &row, getUserQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrUserNotFound
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return row.toCore(), nil
}

func (db *DB) ListUsers(ctx context.Context, role string) ([]core.User, error) {
	var rows []userRow

	if role == "" {
		if err := db.conn.SelectContext(ctx, &rows, getAllUsersQuery); err != nil {
			return nil, fmt.Errorf("error getting users: %w", err)
		}
	} else {
		if err := db.conn.SelectContext(ctx, &rows, getUsersByRoleQuery, role); err != nil {
			return nil, fmt.Errorf("error getting users by role: %w", err)
		}
	}

	users := make([]core.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, *row.toCore())
	}

	return users, nil
}

func (db *DB) UpdateUser(ctx context.Context, input core.UpdateUserInput) (*core.User, error) {
	var row userRow
	if err := db.conn.GetContext(
		ctx,
		&row,
		updateUserQuery,
		input.ID,
		input.Name,
		input.Email,
		input.Password,
		input.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrUserNotFound
		}
		return nil, fmt.Errorf("error updating user: %w", err)
	}
	return row.toCore(), nil
}

func (db *DB) DeleteUser(ctx context.Context, id string) error {
	result, err := db.conn.ExecContext(ctx, deleteUserQuery, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected users: %w", err)
	}

	if affected == 0 {
		return core.ErrUserNotFound
	}

	return nil
}
