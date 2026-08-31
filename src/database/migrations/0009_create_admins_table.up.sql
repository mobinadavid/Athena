create table if not exists admins
(
    id         bigserial    primary key,
    uuid       uuid         default uuid_generate_v4(),
    is_active  boolean      default false,
    first_name varchar(255) default NULL::character varying,
    last_name  varchar(255) default NULL::character varying,
    username varchar(255) UNIQUE default NULL::character varying,
    password   text,
    mobile     varchar(100) not null,
    email      varchar(100) default NULL::character varying,
    totp_secret bytea default NULL,
    totp_secret_url bytea default NULL,
    two_fa_enabled boolean default false,
    recovery_codes text[] default NULL,
    admin_image_uuid      varchar(255),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create unique index if not exists idx_admins_uuid
    on admins (uuid);

create index if not exists idx_admins_deleted_at
    on admins (deleted_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_admins_username
    ON admins(username)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_admins_email
    ON admins(email)
    WHERE deleted_at IS NULL;