-- +goose Up
create table if not exists refreshlist(
    refresh_hash text primary key,
    user_login uuid not null,
    user_agent text,
    expired_at timestamp default now()+'11 DAY'
)

-- +goose Down
drop table if exists refreshlist;
