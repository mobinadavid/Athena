create table if not exists internet_gateway_payments
(
    id bigserial primary key,
    uuid uuid default uuid_generate_v4() unique,
    amount numeric(15, 6) default 0 constraint chk_orders_amount check (amount >= (0):: numeric),
    issuer_reference_number varchar(255) default null unique,
    ipg varchar(255) not null,
    receipt jsonb default null,
    callback_url varchar(255) default null,
    status varchar(255) default 'pending'::character varying not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone
                                                                  );

