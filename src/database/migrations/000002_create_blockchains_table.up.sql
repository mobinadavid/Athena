create table if not exists blockchains
(
    id          bigserial primary key,
    uuid        uuid default uuid_generate_v4(),
    title       jsonb,
    symbol      text not null,
    is_active   boolean default true,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
                              );

create index if not exists idx_blockchains_deleted_at
    on blockchains (deleted_at);
