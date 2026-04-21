create table if not exists projects (
    id uuid primary key default gen_random_uuid(),
    project_code text not null unique,
    name text not null,
    description text not null default '',
    contract_number uuid not null unique,
    status text not null check (status in ('draft', 'planned', 'active', 'paused', 'completed', 'cancelled')) default 'draft',
    customer_name text not null default '',
    manager_id uuid not null,
    planned_start_date date,
    planned_deadline date,
    actual_start_date date,
    actual_deadline date,
    planned_budget numeric(15,2) not null default 0,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_projects_status
    on projects(status);

create index if not exists idx_projects_manager_id
    on projects(manager_id);

create index if not exists idx_projects_created_at
    on projects(created_at desc);
