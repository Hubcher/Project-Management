package postgres

import (
    "context"
    "database/sql"
    "errors"
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
    ID         string       `db:"id"`
    FirstName  string       `db:"first_name"`
    LastName   string       `db:"last_name"`
    MiddleName string       `db:"middle_name"`
    BirthDate  sql.NullTime `db:"birth_date"`
    Phone      string       `db:"phone"`
    Department string       `db:"department"`
    Position   string       `db:"position"`
    AvatarURL  string       `db:"avatar_url"`
    Bio        string       `db:"bio"`
    CreatedAt  time.Time    `db:"created_at"`
    UpdatedAt  time.Time    `db:"updated_at"`
}

func (r userRow) toCore() *core.User {
    var birthDate *time.Time
    if r.BirthDate.Valid {
        day := r.BirthDate.Time.UTC()
        normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
        birthDate = &normalized
    }

    return &core.User{
        ID:         r.ID,
        FirstName:  r.FirstName,
        LastName:   r.LastName,
        MiddleName: r.MiddleName,
        BirthDate:  birthDate,
        Phone:      r.Phone,
        Department: r.Department,
        Position:   r.Position,
        AvatarURL:  r.AvatarURL,
        Bio:        r.Bio,
        CreatedAt:  r.CreatedAt,
        UpdatedAt:  r.UpdatedAt,
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
        insert into users (
            id, first_name, last_name, middle_name, birth_date, phone,
            department, position, avatar_url, bio
        )
        values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        returning id::text, first_name, last_name, middle_name, birth_date, phone,
            department, position, avatar_url, bio, created_at, updated_at;
    `
    getUserQuery = `
        select id::text, first_name, last_name, middle_name, birth_date, phone,
            department, position, avatar_url, bio, created_at, updated_at
        from users
        where id = $1;
    `
    listUsersQuery = `
        select id::text, first_name, last_name, middle_name, birth_date, phone,
            department, position, avatar_url, bio, created_at, updated_at
        from users
        order by created_at desc;
    `
    updateUserQuery = `
        update users
        set first_name = $2,
            last_name = $3,
            middle_name = $4,
            birth_date = $5,
            phone = $6,
            department = $7,
            position = $8,
            avatar_url = $9,
            bio = $10,
            updated_at = current_timestamp
        where id = $1
        returning id::text, first_name, last_name, middle_name, birth_date, phone,
            department, position, avatar_url, bio, created_at, updated_at;
    `
    deleteUserQuery = `
        delete from users
        where id = $1;
    `
)

func nullableDate(value *time.Time) any {
    if value == nil {
        return nil
    }
    return *value
}

func (db *DB) CreateUser(ctx context.Context, input core.CreateUserInput) (*core.User, error) {
    var row userRow
    if err := db.conn.GetContext(
        ctx,
        &row,
        createQuery,
        input.ID,
        input.FirstName,
        input.LastName,
        input.MiddleName,
        nullableDate(input.BirthDate),
        input.Phone,
        input.Department,
        input.Position,
        input.AvatarURL,
        input.Bio,
    ); err != nil {
        return nil, fmt.Errorf("error creating user: %w", err)
    }
    return row.toCore(), nil
}

func (db *DB) GetUser(ctx context.Context, id string) (*core.User, error) {
    var row userRow
    if err := db.conn.GetContext(ctx, &row, getUserQuery, id); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, core.ErrUserNotFound
        }
        return nil, fmt.Errorf("error getting user: %w", err)
    }
    return row.toCore(), nil
}

func (db *DB) ListUsers(ctx context.Context) ([]core.User, error) {
    var rows []userRow
    if err := db.conn.SelectContext(ctx, &rows, listUsersQuery); err != nil {
        return nil, fmt.Errorf("error getting users: %w", err)
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
        input.FirstName,
        input.LastName,
        input.MiddleName,
        nullableDate(input.BirthDate),
        input.Phone,
        input.Department,
        input.Position,
        input.AvatarURL,
        input.Bio,
    ); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
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
