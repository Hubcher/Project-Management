create table if not exists projects (
    contract_number serial primary key,
    name text not null,
    start_date date not null,
    deadline date not null,
    price numeric(12, 2) not null,
    user_id UUID not null,
    created_at timestamp not null default current_timestamp
);

create index if not exists idx_projects_user_id
on projects(user_id)