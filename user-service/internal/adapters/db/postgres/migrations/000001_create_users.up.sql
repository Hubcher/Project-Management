create table if not exists users (
    id uuid primary key,
    name text not null,
    created_at timestamp not null default CURRENT_TIMESTAMP

);
