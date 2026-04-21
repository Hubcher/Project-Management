create table if not exists payment_categories (

    id uuid primary key default gen_random_uuid(),
    name text not null,
    type text not null check (type in ('income', 'expense')),
    parent_id uuid references payment_categories(id) on delete set null,
    is_active boolean not null default true,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    unique(name, type)
);

create index if not exists idx_payment_categories_type
    on payment_categories(type);