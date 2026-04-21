create table if not exists users (
    id uuid primary key,
    first_name text not null,
    last_name text not null,
    middle_name text not null default '',
    birth_date date,
    phone text not null default '',
    department text not null default '',
    position text not null default '',
    avatar_url text not null default '',
    bio text not null default '',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create index if not exists idx_users_created_at
    on users(created_at desc);
