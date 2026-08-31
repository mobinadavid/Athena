create table permission_groups
(
    id    bigserial
        primary key,
    uuid  uuid default uuid_generate_v4(),
    name  varchar(255),
    title jsonb,
    description  varchar(1000),
    is_active       boolean default true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

create unique index idx_permission_groups_name
    on permission_groups (name);

create unique index idx_permission_groups_uuid
    on permission_groups (uuid);

