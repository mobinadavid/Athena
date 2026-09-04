alter table if exists users
    add column if not exists recovery_codes text[] default null;
