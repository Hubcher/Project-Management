create table if not exists calendar_entries (
    id uuid primary key default gen_random_uuid(),
    obligation_id uuid references financial_obligations(id) on delete cascade,
    project_id uuid not null,
    stage_id uuid,
    entry_date date not null,
    type text not null check (type in ('income', 'expense')),
    amount numeric(15,2) not null,
    status text not null check (status in ('planned', 'expected', 'paid', 'overdue', 'cancelled')) default 'planned',
    description text not null default '',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_calendar_entries_project_id
    on calendar_entries(project_id);

create index if not exists idx_calendar_entries_stage_id
    on calendar_entries(stage_id);

create index if not exists idx_calendar_entries_entry_date
    on calendar_entries(entry_date);

create index if not exists idx_calendar_entries_status
    on calendar_entries(status);