create table if not exists payment_requests
(
    id               bigserial primary key,
    uuid             uuid         default uuid_generate_v4(),
    user_id          bigint       not null references users (id),
    blockchain_id    bigint       not null references blockchains (id),
    requested_count  integer      not null default 1,
    expected_amount  numeric,
    status           varchar(50)  not null default 'pending',
    expires_at       timestamp with time zone,
    confirmed_at     timestamp with time zone,
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    deleted_at       timestamp with time zone
);

create unique index if not exists idx_payment_requests_uuid on payment_requests (uuid);
create index if not exists idx_payment_requests_user_id on payment_requests (user_id);
create index if not exists idx_payment_requests_status on payment_requests (status);
create index if not exists idx_payment_requests_deleted_at on payment_requests (deleted_at);

alter table if exists wallet_addresses
    add column if not exists allocated_to_user_id bigint references users (id),
    add column if not exists payment_request_id bigint references payment_requests (id);

create index if not exists idx_wallet_addresses_allocated_to_user_id on wallet_addresses (allocated_to_user_id);
create index if not exists idx_wallet_addresses_payment_request_id on wallet_addresses (payment_request_id);

alter table if exists deposits
    add column if not exists uuid uuid default uuid_generate_v4(),
    add column if not exists user_id bigint references users (id),
    add column if not exists wallet_address_id bigint references wallet_addresses (id),
    add column if not exists payment_request_id bigint references payment_requests (id),
    add column if not exists blockchain_id bigint references blockchains (id),
    add column if not exists from_address text,
    add column if not exists to_address text,
    add column if not exists amount numeric,
    add column if not exists fee numeric,
    add column if not exists confirmations integer default 0,
    add column if not exists status varchar(50) default 'confirmed',
    add column if not exists notified_at timestamp with time zone,
    add column if not exists block_number bigint,
    add column if not exists paid_at timestamp with time zone;

create unique index if not exists idx_deposits_uuid on deposits (uuid);
create unique index if not exists idx_deposits_transaction_hash on deposits (transaction_hash);
create index if not exists idx_deposits_user_id on deposits (user_id);
create index if not exists idx_deposits_payment_request_id on deposits (payment_request_id);

create table if not exists notifications
(
    id         bigserial primary key,
    uuid       uuid         default uuid_generate_v4(),
    user_id    bigint       not null references users (id),
    type       varchar(100) not null,
    title      varchar(255) not null,
    body       text,
    data       jsonb,
    is_read    boolean      default false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create unique index if not exists idx_notifications_uuid on notifications (uuid);
create index if not exists idx_notifications_user_id on notifications (user_id);
create index if not exists idx_notifications_is_read on notifications (is_read);
create index if not exists idx_notifications_deleted_at on notifications (deleted_at);
