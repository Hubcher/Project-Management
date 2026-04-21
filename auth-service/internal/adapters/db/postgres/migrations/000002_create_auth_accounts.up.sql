create extension if not exists pgcrypto;

create table if not exists auth_accounts (
    user_id uuid primary key,
    email text not null unique,
    password_hash text not null,
    role text not null check (role in ('admin', 'user')) default 'user',
    is_active boolean not null default true,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    password_changed_at timestamp not null default current_timestamp
);

create index if not exists idx_auth_accounts_email
    on auth_accounts(email);

