create table if not exists roles
(
    id         bigserial
                       primary key,
    uuid            uuid default uuid_generate_v4(),
    name            varchar(255),
    title           varchar(255),
    description     varchar(1000),
    is_active       boolean default true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


create unique index if not exists idx_roles_name
    on roles (name);

create unique index if not exists idx_roles_uuid
    on roles (uuid);

