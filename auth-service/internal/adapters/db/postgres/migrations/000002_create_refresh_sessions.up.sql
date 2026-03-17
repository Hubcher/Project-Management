create table if not exists refresh_sessions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references auth_accounts(user_id) on delete cascade,
    token_hash text not null unique,
    user_agent text,
    ip inet,
    expires_at timestamp not null,
    created_at timestamp not null default current_timestamp,
    revoked_at timestamp
);

create index if not exists idx_refresh_sessions_user_id
on refresh_sessions(user_id);

create index if not exists idx_refresh_sessions_expires_at
on refresh_sessions(expires_at);