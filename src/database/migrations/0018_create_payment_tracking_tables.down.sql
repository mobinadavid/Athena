drop table if exists notifications;

alter table if exists payment_requests
    drop column if exists confirmed_at;

alter table if exists deposits
    drop column if exists uuid,
    drop column if exists user_id,
    drop column if exists wallet_address_id,
    drop column if exists payment_request_id,
    drop column if exists blockchain_id,
    drop column if exists from_address,
    drop column if exists to_address,
    drop column if exists amount,
    drop column if exists fee,
    drop column if exists confirmations,
    drop column if exists status,
    drop column if exists notified_at,
    drop column if exists block_number,
    drop column if exists paid_at;

alter table if exists wallet_addresses
    drop column if exists allocated_to_user_id,
    drop column if exists payment_request_id;

drop table if exists payment_requests;
