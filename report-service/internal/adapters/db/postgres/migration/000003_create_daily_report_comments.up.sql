create table if not exists daily_report_comments (
    id uuid primary key default gen_random_uuid(),
    report_id uuid not null references daily_reports(id) on delete cascade,
    author_user_id uuid not null,
    comment text not null,
    created_at timestamp not null default current_timestamp
);

create index if not exists idx_daily_report_comments_report_id
    on daily_report_comments(report_id);