create table if not exists access_tokens
(
    id                       bigserial primary key,
    uuid                     uuid,
    owner_id                 bigint,
    owner_type               text,
    access_token             text not null,
    access_token_expires_at  timestamp with time zone,
    refresh_token            text not null,
    ip                       varchar(255) not null,
    user_agent               varchar(255) not null,
    refresh_token_expires_at timestamp with time zone,
    last_used_at             timestamp with time zone,
    created_at               timestamp with time zone,
    updated_at               timestamp with time zone,
    deleted_at               timestamp with time zone
);

create index if not exists idx_access_tokens_deleted_at
    on access_tokens (deleted_at);

create unique index if not exists idx_access_tokens_uuid
    on access_tokens (uuid);

