package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Hubcher/project-management/project-service/internal/core"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

type projectRow struct {
	ID               string     `db:"id"`
	ProjectCode      string     `db:"project_code"`
	Name             string     `db:"name"`
	Description      string     `db:"description"`
	ContractNumber   string     `db:"contract_number"`
	Status           string     `db:"status"`
	CustomerName     string     `db:"customer_name"`
	ManagerID        string     `db:"manager_id"`
	PlannedStartDate *time.Time `db:"planned_start_date"`
	PlannedDeadline  *time.Time `db:"planned_deadline"`
	ActualStartDate  *time.Time `db:"actual_start_date"`
	ActualDeadline   *time.Time `db:"actual_deadline"`
	PlannedBudget    string     `db:"planned_budget"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

type stageRow struct {
	ID               string     `db:"id"`
	ProjectID        string     `db:"project_id"`
	Name             string     `db:"name"`
	Description      string     `db:"description"`
	SequenceNumber   int32      `db:"sequence_number"`
	Status           string     `db:"status"`
	PlannedStartDate *time.Time `db:"planned_start_date"`
	PlannedEndDate   *time.Time `db:"planned_end_date"`
	ActualStartDate  *time.Time `db:"actual_start_date"`
	ActualEndDate    *time.Time `db:"actual_end_date"`
	PlannedIncome    string     `db:"planned_income"`
	PlannedExpense   string     `db:"planned_expense"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

type memberRow struct {
	ID            string     `db:"id"`
	ProjectID     string     `db:"project_id"`
	UserID        string     `db:"user_id"`
	RoleInProject string     `db:"role_in_project"`
	IsActive      bool       `db:"is_active"`
	JoinedAt      time.Time  `db:"joined_at"`
	LeftAt        *time.Time `db:"left_at"`
}

type eventRow struct {
	ID          string     `db:"id"`
	ProjectID   string     `db:"project_id"`
	StageID     string     `db:"stage_id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	PlannedDate *time.Time `db:"planned_date"`
	ActualDate  *time.Time `db:"actual_date"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (r projectRow) toCore() *core.Project {
	return &core.Project{
		ID:               r.ID,
		ProjectCode:      r.ProjectCode,
		Name:             r.Name,
		Description:      r.Description,
		ContractNumber:   r.ContractNumber,
		Status:           core.ProjectStatus(r.Status),
		CustomerName:     r.CustomerName,
		ManagerID:        r.ManagerID,
		PlannedStartDate: r.PlannedStartDate,
		PlannedDeadline:  r.PlannedDeadline,
		ActualStartDate:  r.ActualStartDate,
		ActualDeadline:   r.ActualDeadline,
		PlannedBudget:    r.PlannedBudget,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func (r stageRow) toCore() *core.ProjectStage {
	return &core.ProjectStage{
		ID:               r.ID,
		ProjectID:        r.ProjectID,
		Name:             r.Name,
		Description:      r.Description,
		SequenceNumber:   r.SequenceNumber,
		Status:           core.StageStatus(r.Status),
		PlannedStartDate: r.PlannedStartDate,
		PlannedEndDate:   r.PlannedEndDate,
		ActualStartDate:  r.ActualStartDate,
		ActualEndDate:    r.ActualEndDate,
		PlannedIncome:    r.PlannedIncome,
		PlannedExpense:   r.PlannedExpense,
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}
}

func (r memberRow) toCore() *core.ProjectMember {
	return &core.ProjectMember{
		ID:            r.ID,
		ProjectID:     r.ProjectID,
		UserID:        r.UserID,
		RoleInProject: r.RoleInProject,
		IsActive:      r.IsActive,
		JoinedAt:      r.JoinedAt,
		LeftAt:        r.LeftAt,
	}
}

func (r eventRow) toCore() *core.ProjectEvent {
	return &core.ProjectEvent{
		ID:          r.ID,
		ProjectID:   r.ProjectID,
		StageID:     r.StageID,
		Name:        r.Name,
		Description: r.Description,
		PlannedDate: r.PlannedDate,
		ActualDate:  r.ActualDate,
		Status:      core.EventStatus(r.Status),
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

const (
	projectColumns = `
		id::text,
		project_code,
		name,
		description,
		contract_number::text as contract_number,
		status,
		customer_name,
		manager_id::text as manager_id,
		planned_start_date,
		planned_deadline,
		actual_start_date,
		actual_deadline,
		planned_budget::text as planned_budget,
		created_at,
		updated_at`
	stageColumns = `
		id::text,
		project_id::text as project_id,
		name,
		description,
		sequence_number,
		status,
		planned_start_date,
		planned_end_date,
		actual_start_date,
		actual_end_date,
		planned_income::text as planned_income,
		planned_expense::text as planned_expense,
		created_at,
		updated_at`
	memberColumns = `
		id::text,
		project_id::text as project_id,
		user_id::text as user_id,
		role_in_project,
		is_active,
		joined_at,
		left_at`
	eventColumns = `
		id::text,
		project_id::text as project_id,
		coalesce(stage_id::text, '') as stage_id,
		name,
		description,
		planned_date,
		actual_date,
		status,
		created_at,
		updated_at`
)

func (db *DB) CreateProject(ctx context.Context, input core.CreateProjectInput) (*core.Project, error) {
	query := `
		insert into projects (
			project_code, name, description, contract_number, status, customer_name, manager_id,
			planned_start_date, planned_deadline, actual_start_date, actual_deadline, planned_budget
		)
		values ($1, $2, $3, $4::uuid, $5, $6, $7::uuid, $8, $9, $10, $11, $12)
		returning ` + projectColumns + `;`
	var row projectRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ProjectCode, input.Name, input.Description, input.ContractNumber, string(input.Status), input.CustomerName, input.ManagerID,
		input.PlannedStartDate, input.PlannedDeadline, input.ActualStartDate, input.ActualDeadline, input.PlannedBudget,
	); err != nil {
		return nil, mapWriteErr(err, core.ErrInvalidProject)
	}
	return row.toCore(), nil
}

func (db *DB) GetProject(ctx context.Context, id string) (*core.Project, error) {
	query := `select ` + projectColumns + ` from projects where id = $1::uuid;`
	var row projectRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrProjectNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListProjects(ctx context.Context, participantUserID string) ([]core.Project, error) {
	query := `
		select ` + projectColumns + `
		from projects p
		where (
			$1 = '' or p.manager_id::text = $1 or exists (
				select 1
				from project_members pm
				where pm.project_id = p.id and pm.user_id::text = $1 and pm.is_active = true
			)
		)
		order by p.created_at desc;`
	var rows []projectRow
	if err := db.conn.SelectContext(ctx, &rows, query, participantUserID); err != nil {
		return nil, err
	}
	projects := make([]core.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, *row.toCore())
	}
	return projects, nil
}

func (db *DB) UpdateProject(ctx context.Context, input core.UpdateProjectInput) (*core.Project, error) {
	query := `
		update projects
		set project_code = $2,
			name = $3,
			description = $4,
			contract_number = $5::uuid,
			status = $6,
			customer_name = $7,
			manager_id = $8::uuid,
			planned_start_date = $9,
			planned_deadline = $10,
			actual_start_date = $11,
			actual_deadline = $12,
			planned_budget = $13,
			updated_at = current_timestamp
		where id = $1::uuid
		returning ` + projectColumns + `;`
	var row projectRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ID, input.ProjectCode, input.Name, input.Description, input.ContractNumber, string(input.Status), input.CustomerName, input.ManagerID,
		input.PlannedStartDate, input.PlannedDeadline, input.ActualStartDate, input.ActualDeadline, input.PlannedBudget,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrProjectNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidProject)
	}
	return row.toCore(), nil
}

func (db *DB) DeleteProject(ctx context.Context, id string) error {
	return deleteByID(ctx, db.conn, `delete from projects where id = $1::uuid;`, id, core.ErrProjectNotFound)
}

func (db *DB) CreateStage(ctx context.Context, input core.CreateProjectStageInput) (*core.ProjectStage, error) {
	query := `
		insert into project_stages (
			project_id, name, description, sequence_number, status,
			planned_start_date, planned_end_date, actual_start_date, actual_end_date, planned_income, planned_expense
		)
		values ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		returning ` + stageColumns + `;`
	var row stageRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ProjectID, input.Name, input.Description, input.SequenceNumber, string(input.Status),
		input.PlannedStartDate, input.PlannedEndDate, input.ActualStartDate, input.ActualEndDate, input.PlannedIncome, input.PlannedExpense,
	); err != nil {
		return nil, mapWriteErr(err, core.ErrInvalidStage)
	}
	return row.toCore(), nil
}

func (db *DB) GetStage(ctx context.Context, id string) (*core.ProjectStage, error) {
	query := `select ` + stageColumns + ` from project_stages where id = $1::uuid;`
	var row stageRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrStageNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListStages(ctx context.Context, projectID string) ([]core.ProjectStage, error) {
	query := `select ` + stageColumns + ` from project_stages where project_id = $1::uuid order by sequence_number asc, created_at asc;`
	var rows []stageRow
	if err := db.conn.SelectContext(ctx, &rows, query, projectID); err != nil {
		return nil, err
	}
	stages := make([]core.ProjectStage, 0, len(rows))
	for _, row := range rows {
		stages = append(stages, *row.toCore())
	}
	return stages, nil
}

func (db *DB) UpdateStage(ctx context.Context, input core.UpdateProjectStageInput) (*core.ProjectStage, error) {
	query := `
		update project_stages
		set name = $2,
			description = $3,
			sequence_number = $4,
			status = $5,
			planned_start_date = $6,
			planned_end_date = $7,
			actual_start_date = $8,
			actual_end_date = $9,
			planned_income = $10,
			planned_expense = $11,
			updated_at = current_timestamp
		where id = $1::uuid
		returning ` + stageColumns + `;`
	var row stageRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ID, input.Name, input.Description, input.SequenceNumber, string(input.Status),
		input.PlannedStartDate, input.PlannedEndDate, input.ActualStartDate, input.ActualEndDate, input.PlannedIncome, input.PlannedExpense,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrStageNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidStage)
	}
	return row.toCore(), nil
}

func (db *DB) DeleteStage(ctx context.Context, id string) error {
	return deleteByID(ctx, db.conn, `delete from project_stages where id = $1::uuid;`, id, core.ErrStageNotFound)
}

func (db *DB) CreateMember(ctx context.Context, input core.CreateProjectMemberInput) (*core.ProjectMember, error) {
	query := `
		insert into project_members (project_id, user_id, role_in_project)
		values ($1::uuid, $2::uuid, $3)
		returning ` + memberColumns + `;`
	var row memberRow
	if err := db.conn.GetContext(ctx, &row, query, input.ProjectID, input.UserID, input.RoleInProject); err != nil {
		return nil, mapWriteErr(err, core.ErrInvalidMember)
	}
	return row.toCore(), nil
}

func (db *DB) GetMember(ctx context.Context, id string) (*core.ProjectMember, error) {
	query := `select ` + memberColumns + ` from project_members where id = $1::uuid;`
	var row memberRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrMemberNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListMembers(ctx context.Context, projectID string) ([]core.ProjectMember, error) {
	query := `select ` + memberColumns + ` from project_members where project_id = $1::uuid order by joined_at asc;`
	var rows []memberRow
	if err := db.conn.SelectContext(ctx, &rows, query, projectID); err != nil {
		return nil, err
	}
	members := make([]core.ProjectMember, 0, len(rows))
	for _, row := range rows {
		members = append(members, *row.toCore())
	}
	return members, nil
}

func (db *DB) UpdateMember(ctx context.Context, input core.UpdateProjectMemberInput) (*core.ProjectMember, error) {
	query := `
		update project_members
		set role_in_project = $2,
			is_active = $3,
			left_at = $4
		where id = $1::uuid
		returning ` + memberColumns + `;`
	var row memberRow
	if err := db.conn.GetContext(ctx, &row, query, input.ID, input.RoleInProject, input.IsActive, input.LeftAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrMemberNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidMember)
	}
	return row.toCore(), nil
}

func (db *DB) DeleteMember(ctx context.Context, id string) error {
	return deleteByID(ctx, db.conn, `delete from project_members where id = $1::uuid;`, id, core.ErrMemberNotFound)
}

func (db *DB) CreateEvent(ctx context.Context, input core.CreateProjectEventInput) (*core.ProjectEvent, error) {
	query := `
		insert into project_events (project_id, stage_id, name, description, planned_date, actual_date, status)
		values ($1::uuid, $2::uuid, $3, $4, $5, $6, $7)
		returning ` + eventColumns + `;`
	var row eventRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ProjectID, nullableUUID(input.StageID), input.Name, input.Description, input.PlannedDate, input.ActualDate, string(input.Status),
	); err != nil {
		return nil, mapWriteErr(err, core.ErrInvalidEvent)
	}
	return row.toCore(), nil
}

func (db *DB) GetEvent(ctx context.Context, id string) (*core.ProjectEvent, error) {
	query := `select ` + eventColumns + ` from project_events where id = $1::uuid;`
	var row eventRow
	if err := db.conn.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrEventNotFound
		}
		return nil, err
	}
	return row.toCore(), nil
}

func (db *DB) ListEvents(ctx context.Context, projectID string) ([]core.ProjectEvent, error) {
	query := `select ` + eventColumns + ` from project_events where project_id = $1::uuid order by created_at desc;`
	var rows []eventRow
	if err := db.conn.SelectContext(ctx, &rows, query, projectID); err != nil {
		return nil, err
	}
	events := make([]core.ProjectEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, *row.toCore())
	}
	return events, nil
}

func (db *DB) UpdateEvent(ctx context.Context, input core.UpdateProjectEventInput) (*core.ProjectEvent, error) {
	query := `
		update project_events
		set stage_id = $2::uuid,
			name = $3,
			description = $4,
			planned_date = $5,
			actual_date = $6,
			status = $7,
			updated_at = current_timestamp
		where id = $1::uuid
		returning ` + eventColumns + `;`
	var row eventRow
	if err := db.conn.GetContext(ctx, &row, query,
		input.ID, nullableUUID(input.StageID), input.Name, input.Description, input.PlannedDate, input.ActualDate, string(input.Status),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrEventNotFound
		}
		return nil, mapWriteErr(err, core.ErrInvalidEvent)
	}
	return row.toCore(), nil
}

func (db *DB) DeleteEvent(ctx context.Context, id string) error {
	return deleteByID(ctx, db.conn, `delete from project_events where id = $1::uuid;`, id, core.ErrEventNotFound)
}

func deleteByID(ctx context.Context, conn *sqlx.DB, query, id string, notFound error) error {
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
	case strings.Contains(message, "foreign key constraint") || strings.Contains(message, "invalid input syntax"):
		return invalidErr
	default:
		return fmt.Errorf("db write: %w", err)
	}
}
