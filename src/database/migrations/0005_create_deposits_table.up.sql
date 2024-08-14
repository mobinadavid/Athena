create table if not exists deposits (

    id                bigserial primary key,
    transaction_hash  text,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone
                                    );

create index if not exists idx_deposits_deleted_at
    on deposits (deleted_at);

