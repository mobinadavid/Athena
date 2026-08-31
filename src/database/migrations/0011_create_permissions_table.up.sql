create table permissions
(
    id    bigserial
        primary key,
    uuid  uuid default uuid_generate_v4(),
    name  varchar(255),
    title jsonb,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

create unique index idx_permissions_name
    on permissions (name);

create unique index idx_permissions_uuid
    on permissions (uuid);

