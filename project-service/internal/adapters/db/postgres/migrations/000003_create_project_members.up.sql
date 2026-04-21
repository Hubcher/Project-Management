create table if not exists project_members (
    id uuid primary key default gen_random_uuid(),
    project_id uuid not null references projects(id) on delete cascade,
    user_id uuid not null,
    role_in_project text not null default 'member',
    is_active boolean not null default true,
    joined_at timestamp not null default current_timestamp,
    left_at timestamp,
    unique(project_id, user_id)
);

create index if not exists idx_project_members_project_id
    on project_members(project_id);

create index if not exists idx_project_members_user_id
    on project_members(user_id);