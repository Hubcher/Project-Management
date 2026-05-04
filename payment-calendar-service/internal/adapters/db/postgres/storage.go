package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Hubcher/project-management/payment-calendar-service/internal/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

type paymentRow struct {
	ID          string     `db:"id"`
	ProjectID   string     `db:"project_id"`
	StageID     string     `db:"stage_id"`
	Type        string     `db:"type"`
	Status      string     `db:"status"`
	Amount      string     `db:"amount"`
	Currency    string     `db:"currency"`
	PlannedDate time.Time  `db:"planned_date"`
	ActualDate  *time.Time `db:"actual_date"`
	Description string     `db:"description"`
	CreatedBy   string     `db:"created_by"`
	PaidBy      string     `db:"paid_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (r paymentRow) toCore() *core.Payment {
	return &core.Payment{
		ID:          r.ID,
		ProjectID:   r.ProjectID,
		StageID:     r.StageID,
		Type:        core.PaymentType(r.Type),
		Status:      core.PaymentStatus(r.Status),
		Amount:      r.Amount,
		Currency:    r.Currency,
		PlannedDate: r.PlannedDate,
		ActualDate:  r.ActualDate,
		Description: r.Description,
		CreatedBy:   r.CreatedBy,
		PaidBy:      r.PaidBy,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
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

const paymentColumns = `
	id::text,
	project_id::text as project_id,
	coalesce(stage_id::text, '') as stage_id,
	type,
	status,
	amount::text as amount,
	currency,
	planned_date,
	actual_date,
	description,
	coalesce(created_by::text, '') as created_by,
	coalesce(paid_by::text, '') as paid_by,
	created_at,
	updated_at`

func (db *DB) CreatePayment(ctx context.Context, input core.CreatePaymentInput) (*core.Payment, error) {
	query := `
		insert into payments (
			project_id, stage_id, type, status, amount, currency, planned_date, description, created_by
		)
		values ($1::uuid, $2::uuid, $3, 'planned', $4, $5, $6, $7, $8::uuid)
		returning ` + paymentColumns + `;`
	var row paymentRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ProjectID,
		nullableUUID(input.StageID),
		string(input.Type),
		input.Amount,
		input.Currency,
		input.PlannedDate,
		input.Description,
		nullableUUID(input.CreatedBy),
	); err != nil {
		return nil, mapWriteErr(err, core.ErrInvalidPayment)
	}
	return row.toCore(), nil
}

func (db *DB) GetPayment(ctx context.Context, id string) (*core.Payment, error) {
	query := `select ` + paymentColumns + ` from payments where id = $1::uuid;`
	var row paymentRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrPaymentNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListPayments(ctx context.Context, filter core.ListPaymentsFilter) ([]core.Payment, error) {
	query := `
		select ` + paymentColumns + `
		from payments
		where project_id = $1::uuid
		  and ($2 = '' or coalesce(stage_id::text, '') = $2)
		  and ($3 = '' or type = $3)
		  and ($4 = '' or status = $4)
		  and ($5::date is null or planned_date >= $5::date)
		  and ($6::date is null or planned_date <= $6::date)
		order by planned_date asc, created_at asc;`
	var rows []paymentRow
	if err := db.conn.SelectContext(ctx, &rows, query,
		filter.ProjectID,
		filter.StageID,
		string(filter.Type),
		string(filter.Status),
		filter.DateFrom,
		filter.DateTo,
	); err != nil {
		return nil, err
	}
	payments := make([]core.Payment, 0, len(rows))
	for _, row := range rows {
		payments = append(payments, *row.toCore())
	}
	return payments, nil
}

func (db *DB) UpdatePayment(ctx context.Context, input core.UpdatePaymentInput) (*core.Payment, error) {
	query := `
		update payments
		set stage_id = $2::uuid,
			type = $3,
			status = $4,
			amount = $5,
			currency = $6,
			planned_date = $7,
			actual_date = $8,
			description = $9,
			paid_by = $10::uuid,
			updated_at = current_timestamp
		where id = $1::uuid
		returning ` + paymentColumns + `;`
	var row paymentRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ID,
		nullableUUID(input.StageID),
		string(input.Type),
		string(input.Status),
		input.Amount,
		input.Currency,
		input.PlannedDate,
		input.ActualDate,
		input.Description,
		nullableUUID(input.PaidBy),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrPaymentNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidPayment)
	}
	return row.toCore(), nil
}

func (db *DB) DeletePayment(ctx context.Context, id string) error {
	result, err := db.conn.ExecContext(ctx, `delete from payments where id = $1::uuid;`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return core.ErrPaymentNotFound
	}
	return nil
}

func (db *DB) MarkPaymentPaid(ctx context.Context, input core.MarkPaymentPaidInput) (*core.Payment, error) {
	query := `
		update payments
		set status = 'paid',
			actual_date = $2,
			paid_by = $3::uuid,
			updated_at = current_timestamp
		where id = $1::uuid
		returning ` + paymentColumns + `;`
	var row paymentRow
	if err := db.conn.GetContext(ctx, &row, query, input.ID, input.ActualDate, nullableUUID(input.PaidBy)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrPaymentNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidPayment)
	}
	return row.toCore(), nil
}

func nullableUUID(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

func mapWriteErr(err error, invalidErr error) error {
	if err == nil {
		return nil
	}
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "duplicate key") || strings.Contains(message, "unique constraint"):
		return core.ErrAlreadyExists
	case strings.Contains(message, "foreign key constraint") || strings.Contains(message, "invalid input syntax") || strings.Contains(message, "check constraint"):
		return invalidErr
	default:
		return fmt.Errorf("db write: %w", err)
	}
}
