create table if not exists financial_obligations (

    id uuid primary key default gen_random_uuid(),
    project_id uuid not null,
    stage_id uuid,
    counterparty_id uuid references counterparties(id) on delete set null,
    category_id uuid not null references payment_categories(id) on delete restrict,
    type text not null check (type in ('income', 'expense')),
    basis text not null default '',
    description text not null default '',
    amount numeric(15,2) not null,
    currency text not null default 'RUB',
    due_date date not null,
    status text not null check (status in ('planned', 'approved', 'partially_paid', 'paid', 'overdue', 'cancelled')) default 'planned',
    responsible_user_id uuid,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_financial_obligations_project_id
    on financial_obligations(project_id);

create index if not exists idx_financial_obligations_stage_id
    on financial_obligations(stage_id);

create index if not exists idx_financial_obligations_due_date
    on financial_obligations(due_date);

create index if not exists idx_financial_obligations_status
    on financial_obligations(status);