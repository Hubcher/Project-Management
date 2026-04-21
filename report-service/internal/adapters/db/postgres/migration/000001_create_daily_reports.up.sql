create extension if not exists pgcrypto;

create table if not exists daily_reports (

    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    report_date date not null,
    status text not null check (status in ('draft', 'submitted', 'approved', 'rejected')) default 'draft',
    total_hours numeric(6,2) not null default 0,
    summary text not null default '',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    unique(user_id, report_date)
);

create index if not exists idx_daily_reports_user_id
    on daily_reports(user_id);

create index if not exists idx_daily_reports_report_date
    on daily_reports(report_date desc);

create index if not exists idx_daily_reports_status
    on daily_reports(status);