package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Hubcher/project-management/report-service/internal/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

type reportRow struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	ReportDate time.Time `db:"report_date"`
	Status     string    `db:"status"`
	TotalHours string    `db:"total_hours"`
	Summary    string    `db:"summary"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type entryRow struct {
	ID          string    `db:"id"`
	ReportID    string    `db:"report_id"`
	ProjectID   string    `db:"project_id"`
	StageID     string    `db:"stage_id"`
	WorkType    string    `db:"work_type"`
	Description string    `db:"description"`
	HoursSpent  string    `db:"hours_spent"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type commentRow struct {
	ID           string    `db:"id"`
	ReportID     string    `db:"report_id"`
	AuthorUserID string    `db:"author_user_id"`
	Comment      string    `db:"comment"`
	CreatedAt    time.Time `db:"created_at"`
}

func (r reportRow) toCore() *core.DailyReport {
	return &core.DailyReport{
		ID:         r.ID,
		UserID:     r.UserID,
		ReportDate: r.ReportDate,
		Status:     core.ReportStatus(r.Status),
		TotalHours: r.TotalHours,
		Summary:    r.Summary,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func (r entryRow) toCore() *core.DailyReportEntry {
	return &core.DailyReportEntry{
		ID:          r.ID,
		ReportID:    r.ReportID,
		ProjectID:   r.ProjectID,
		StageID:     r.StageID,
		WorkType:    r.WorkType,
		Description: r.Description,
		HoursSpent:  r.HoursSpent,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func (r commentRow) toCore() *core.DailyReportComment {
	return &core.DailyReportComment{
		ID:           r.ID,
		ReportID:     r.ReportID,
		AuthorUserID: r.AuthorUserID,
		Comment:      r.Comment,
		CreatedAt:    r.CreatedAt,
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
	reportColumns = `
		id::text,
		user_id::text as user_id,
		report_date,
		status,
		total_hours::text as total_hours,
		summary,
		created_at,
		updated_at`
	entryColumns = `
		id::text,
		report_id::text as report_id,
		project_id::text as project_id,
		coalesce(stage_id::text, '') as stage_id,
		work_type,
		description,
		hours_spent::text as hours_spent,
		created_at,
		updated_at`
	commentColumns = `
		id::text,
		report_id::text as report_id,
		author_user_id::text as author_user_id,
		comment,
		created_at`
)

func (db *DB) CreateReport(ctx context.Context, input core.CreateDailyReportInput) (*core.DailyReport, error) {
	query := `
		insert into daily_reports (user_id, report_date, status, summary)
		values ($1::uuid, $2, $3, $4)
		returning ` + reportColumns + `;`
	var row reportRow
	if err := db.conn.GetContext(ctx, &row, query, input.UserID, input.ReportDate, string(input.Status), input.Summary); err != nil {
		return nil, mapWriteErr(err, core.ErrInvalidReport)
	}
	return row.toCore(), nil
}

func (db *DB) GetReport(ctx context.Context, id string) (*core.DailyReport, error) {
	query := `select ` + reportColumns + ` from daily_reports where id = $1::uuid;`
	var row reportRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrReportNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListReports(ctx context.Context, filter core.ListDailyReportsFilter) ([]core.DailyReport, error) {
	query := `
		select ` + reportColumns + `
		from daily_reports
		where ($1 = '' or user_id::text = $1)
		  and ($2 = '' or status = $2)
		  and ($3::date is null or report_date >= $3::date)
		  and ($4::date is null or report_date <= $4::date)
		order by report_date desc, created_at desc;`
	var rows []reportRow
	if err := db.conn.SelectContext(ctx, &rows, query, filter.UserID, string(filter.Status), filter.DateFrom, filter.DateTo); err != nil {
		return nil, err
	}
	reports := make([]core.DailyReport, 0, len(rows))
	for _, row := range rows {
		reports = append(reports, *row.toCore())
	}
	return reports, nil
}

func (db *DB) UpdateReport(ctx context.Context, input core.UpdateDailyReportInput) (*core.DailyReport, error) {
	query := `
		update daily_reports
		set status = $2,
			summary = $3,
			updated_at = current_timestamp
		where id = $1::uuid
		returning ` + reportColumns + `;`
	var row reportRow
	if err := db.conn.GetContext(ctx, &row, query, input.ID, string(input.Status), input.Summary); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrReportNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidReport)
	}
	return row.toCore(), nil
}

func (db *DB) DeleteReport(ctx context.Context, id string) error {
	return deleteByID(ctx, db.conn, `delete from daily_reports where id = $1::uuid;`, id, core.ErrReportNotFound)
}

func (db *DB) CreateEntry(ctx context.Context, input core.CreateDailyReportEntryInput) (*core.DailyReportEntry, error) {
	var row entryRow
	err := db.withTx(ctx, func(tx *sqlx.Tx) error {
		query := `
			insert into daily_report_entries (report_id, project_id, stage_id, work_type, description, hours_spent)
			values ($1::uuid, $2::uuid, $3::uuid, $4, $5, $6)
			returning ` + entryColumns + `;`
		if err := tx.GetContext(ctx, &row, query, input.ReportID, input.ProjectID, nullableUUID(input.StageID), input.WorkType, input.Description, input.HoursSpent); err != nil {
			return mapWriteErr(err, core.ErrInvalidEntry)
		}
		return recalculateReportTotalHours(ctx, tx, input.ReportID)
	})
	if err != nil {
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) GetEntry(ctx context.Context, id string) (*core.DailyReportEntry, error) {
	query := `select ` + entryColumns + ` from daily_report_entries where id = $1::uuid;`
	var row entryRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrEntryNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListEntries(ctx context.Context, reportID string) ([]core.DailyReportEntry, error) {
	query := `select ` + entryColumns + ` from daily_report_entries where report_id = $1::uuid order by created_at asc;`
	var rows []entryRow
	if err := db.conn.SelectContext(ctx, &rows, query, reportID); err != nil {
		return nil, err
	}
	entries := make([]core.DailyReportEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, *row.toCore())
	}
	return entries, nil
}

func (db *DB) UpdateEntry(ctx context.Context, input core.UpdateDailyReportEntryInput) (*core.DailyReportEntry, error) {
	var row entryRow
	err := db.withTx(ctx, func(tx *sqlx.Tx) error {
		query := `
			update daily_report_entries
			set project_id = $2::uuid,
				stage_id = $3::uuid,
				work_type = $4,
				description = $5,
				hours_spent = $6,
				updated_at = current_timestamp
			where id = $1::uuid
			returning ` + entryColumns + `;`
		if err := tx.GetContext(ctx, &row, query, input.ID, input.ProjectID, nullableUUID(input.StageID), input.WorkType, input.Description, input.HoursSpent); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return core.ErrEntryNotFound
			}
			return mapWriteErr(err, core.ErrInvalidEntry)
		}
		return recalculateReportTotalHours(ctx, tx, row.ReportID)
	})
	if err != nil {
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) DeleteEntry(ctx context.Context, id string) error {
	return db.withTx(ctx, func(tx *sqlx.Tx) error {
		reportID, err := getEntryReportID(ctx, tx, id)
		if err != nil {
			return err
		}
		if err = deleteByID(ctx, tx, `delete from daily_report_entries where id = $1::uuid;`, id, core.ErrEntryNotFound); err != nil {
			return err
		}
		return recalculateReportTotalHours(ctx, tx, reportID)
	})
}

func (db *DB) CreateComment(ctx context.Context, input core.CreateDailyReportCommentInput) (*core.DailyReportComment, error) {
	query := `
		insert into daily_report_comments (report_id, author_user_id, comment)
		values ($1::uuid, $2::uuid, $3)
		returning ` + commentColumns + `;`
	var row commentRow
	if err := db.conn.GetContext(ctx, &row, query, input.ReportID, input.AuthorUserID, input.Comment); err != nil {
		return nil, mapWriteErr(err, core.ErrInvalidComment)
	}
	return row.toCore(), nil
}

func (db *DB) GetComment(ctx context.Context, id string) (*core.DailyReportComment, error) {
	query := `select ` + commentColumns + ` from daily_report_comments where id = $1::uuid;`
	var row commentRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrCommentNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListComments(ctx context.Context, reportID string) ([]core.DailyReportComment, error) {
	query := `select ` + commentColumns + ` from daily_report_comments where report_id = $1::uuid order by created_at asc;`
	var rows []commentRow
	if err := db.conn.SelectContext(ctx, &rows, query, reportID); err != nil {
		return nil, err
	}
	comments := make([]core.DailyReportComment, 0, len(rows))
	for _, row := range rows {
		comments = append(comments, *row.toCore())
	}
	return comments, nil
}

func (db *DB) UpdateComment(ctx context.Context, input core.UpdateDailyReportCommentInput) (*core.DailyReportComment, error) {
	query := `
		update daily_report_comments
		set comment = $2
		where id = $1::uuid
		returning ` + commentColumns + `;`
	var row commentRow
	if err := db.conn.GetContext(ctx, &row, query, input.ID, input.Comment); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrCommentNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidComment)
	}
	return row.toCore(), nil
}

func (db *DB) DeleteComment(ctx context.Context, id string) error {
	return deleteByID(ctx, db.conn, `delete from daily_report_comments where id = $1::uuid;`, id, core.ErrCommentNotFound)
}

func (db *DB) withTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := db.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func recalculateReportTotalHours(ctx context.Context, tx *sqlx.Tx, reportID string) error {
	query := `
		update daily_reports
		set total_hours = coalesce((
			select sum(hours_spent)
			from daily_report_entries
			where report_id = $1::uuid
		), 0),
			updated_at = current_timestamp
		where id = $1::uuid;`
	result, err := tx.ExecContext(ctx, query, reportID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return core.ErrReportNotFound
	}
	return nil
}

func getEntryReportID(ctx context.Context, tx *sqlx.Tx, id string) (string, error) {
	var reportID string
	query := `select report_id::text from daily_report_entries where id = $1::uuid;`
	if err := tx.GetContext(ctx, &reportID, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", core.ErrEntryNotFound
		}
		return "", err
	}
	return reportID, nil
}

type execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func deleteByID(ctx context.Context, conn execer, query, id string, notFound error) error {
	result, err := conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return notFound
	}
	return nil
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