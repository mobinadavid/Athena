create table if not exists users
(
    id                           bigserial    primary key,
    uuid                         uuid         default uuid_generate_v4(),
    is_active                    boolean      default true,
    first_name                   varchar(255) default NULL::character varying,
    last_name                    varchar(255) default NULL::character varying,
    password                     text,
    national_identity_code       varchar(255) default NULL::character varying,
    mobile                       varchar(100) not null,
    email                        varchar(100) default NULL::character varying,
    profile_image_uuid           varchar(100) default NULL::character varying,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone DEFAULT NULL
);

create unique index if not exists idx_users_uuid
    on users (uuid);

create index if not exists idx_users_deleted_at
    on users (deleted_at);

create unique index if not exists idx_users_national_identity_code
    on users (national_identity_code);