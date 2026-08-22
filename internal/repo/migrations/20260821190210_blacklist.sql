-- +goose Up
create table if not exists blacklist_refresh (
    refresh_hash text primary key
    created_at timestamp default now()
)


-- +goose Down
drop table if exists blacklist_refresh;