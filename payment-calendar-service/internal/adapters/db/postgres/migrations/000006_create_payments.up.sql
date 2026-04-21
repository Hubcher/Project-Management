create table if not exists payments (
    id uuid primary key default gen_random_uuid(),
    obligation_id uuid references financial_obligations(id) on delete set null,
    calendar_entry_id uuid references calendar_entries(id) on delete set null,
    project_id uuid not null,
    stage_id uuid,
    counterparty_id uuid references counterparties(id) on delete set null,
    type text not null check (type in ('income', 'expense')),
    amount numeric(15,2) not null,
    currency text not null default 'RUB',
    payment_date date not null,
    payment_method text not null default '',
    document_number text not null default '',
    description text not null default '',
    created_by uuid,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_payments_project_id
    on payments(project_id);

create index if not exists idx_payments_stage_id
    on payments(stage_id);

create index if not exists idx_payments_payment_date
    on payments(payment_date desc);

create index if not exists idx_payments_type
    on payments(type);