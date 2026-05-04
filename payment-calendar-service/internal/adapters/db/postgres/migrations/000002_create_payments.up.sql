create table if not exists payments (
    id uuid primary key default gen_random_uuid(),
    project_id uuid not null,
    stage_id uuid,
    type text not null check (type in ('income', 'expense')),
    status text not null check (status in ('planned', 'paid', 'cancelled')) default 'planned',
    amount numeric(15,2) not null check (amount > 0),
    currency text not null default 'RUB',
    planned_date date not null,
    actual_date date,
    description text not null default '',
    created_by uuid,
    paid_by uuid,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    check (
        (status = 'paid' and actual_date is not null)
        or (status <> 'paid' and actual_date is null and paid_by is null)
    )
);

create index if not exists idx_payments_project_id
    on payments(project_id);

create index if not exists idx_payments_stage_id
    on payments(stage_id);

create index if not exists idx_payments_planned_date
    on payments(planned_date);

create index if not exists idx_payments_status
    on payments(status);

create index if not exists idx_payments_type
    on payments(type);
