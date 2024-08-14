create table if not exists wallet_addresses
(
    id               bigserial primary key,
    uuid             uuid default uuid_generate_v4(),
    name             varchar(255) not null,
    wallet_address   varchar(255) not null,
    webhook_url      varchar(255) not null,
    is_active        boolean default true,
    allocated_at     timestamp with time zone,
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    deleted_at       timestamp with time zone,
    blockchain_id    bigserial,
    CONSTRAINT fk_blockchain FOREIGN KEY (blockchain_id) REFERENCES blockchains (id)
                              );

create index if not exists idx_wallet_addresses_deleted_at
    on wallet_addresses (deleted_at);
CREATE INDEX IF NOT EXISTS idx_wallet_addresses_blockchain_id
    ON wallet_addresses (blockchain_id);