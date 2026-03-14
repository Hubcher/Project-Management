create table if not exists users (
    id uuid primary key,
    name text not null,
    email text not null unique,
    password text not null,
    role text not null check (role in ('admin', 'user')) default 'user',
    created_at timestamp not null default CURRENT_TIMESTAMP

);

create index if not exists idx_users_email
on users(email);