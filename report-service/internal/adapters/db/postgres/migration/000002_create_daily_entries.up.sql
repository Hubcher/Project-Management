create table if not exists daily_report_entries (
    id uuid primary key default gen_random_uuid(),
    report_id uuid not null references daily_reports(id) on delete cascade,
    project_id uuid not null,
    stage_id uuid,
    work_type text not null default '',
    description text not null,
    hours_spent numeric(6,2) not null default 0,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_daily_report_entries_report_id
    on daily_report_entries(report_id);

create index if not exists idx_daily_report_entries_project_id
    on daily_report_entries(project_id);

create index if not exists idx_daily_report_entries_stage_id
    on daily_report_entries(stage_id);