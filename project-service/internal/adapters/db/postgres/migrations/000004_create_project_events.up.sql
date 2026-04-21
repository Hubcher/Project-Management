create table if not exists project_events (

    id uuid primary key default gen_random_uuid(),
    project_id uuid not null references projects(id) on delete cascade,
    stage_id uuid references project_stages(id) on delete set null,
    name text not null,
    description text not null default '',
    planned_date date,
    actual_date date,
    status text not null check (status in ('planned', 'reached', 'cancelled')) default 'planned',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_project_events_project_id
    on project_events(project_id);

create index if not exists idx_project_events_stage_id
    on project_events(stage_id);