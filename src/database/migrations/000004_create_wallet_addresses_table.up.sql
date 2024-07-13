create table if not exists wallet_addresses
(
    id          bigserial primary key,
    uuid        uuid default uuid_generate_v4(),
    walletAddress           text not null,
    is_active
    allocated_at timestamp with time zone,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
                              );

create index if not exists idx_wallet_addresses_deleted_at
    on wallet_addresses (deleted_at);
