create table if not exists blockchain_explorers
(
    id              bigserial primary key,
    uuid            uuid default uuid_generate_v4(),
    base_url        text not null,
    is_active       boolean default true,
    is_deafult      text not null,
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone
                                  );

create index if not exists idx_blockchain_explorers_deleted_at
    on blockchain_explorers (deleted_at);
