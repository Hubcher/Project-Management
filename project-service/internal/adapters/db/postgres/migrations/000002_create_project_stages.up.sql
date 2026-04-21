create table if not exists project_stages (
    id uuid primary key default gen_random_uuid(),
    project_id uuid not null references projects(id) on delete cascade,
    name text not null,
    description text not null default '',
    sequence_number integer not null,
    status text not null check (status in ('draft', 'planned', 'in_progress', 'completed', 'cancelled')) default 'draft',
    planned_start_date date,
    planned_end_date date,
    actual_start_date date,
    actual_end_date date,
    planned_income numeric(15,2) not null default 0,
    planned_expense numeric(15,2) not null default 0,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    unique(project_id, sequence_number)
);

create index if not exists idx_project_stages_project_id
    on project_stages(project_id);

create index if not exists idx_project_stages_status
    on project_stages(status);