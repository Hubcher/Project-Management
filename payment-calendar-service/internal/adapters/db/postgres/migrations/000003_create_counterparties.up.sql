create table if not exists counterparties (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    type text not null check (type in ('customer', 'supplier', 'contractor', 'other')),
    contact_info text not null default '',
    tax_id text not null default '',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_counterparties_name
    on counterparties(name);