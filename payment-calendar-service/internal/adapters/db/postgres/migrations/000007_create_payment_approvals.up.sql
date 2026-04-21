create table if not exists payment_approvals (
    id uuid primary key default gen_random_uuid(),
    obligation_id uuid not null references financial_obligations(id) on delete cascade,
    approver_user_id uuid not null,
    status text not null check (status in ('pending', 'approved', 'rejected')) default 'pending',
    comment text not null default '',
    decided_at timestamp,
    created_at timestamp not null default current_timestamp
);

create index if not exists idx_payment_approvals_obligation_id
    on payment_approvals(obligation_id);

create index if not exists idx_payment_approvals_approver_user_id
    on payment_approvals(approver_user_id);