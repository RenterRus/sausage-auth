-- +goose Up
create table if not exists users (
  	user_login text primary key,
    uuid uuid not null default gen_random_uuid(),
    otp_hash text not null,
    otp_link text not null,
    
    confirmed boolean default false,    
    created_at timestamp default now(),
    last_sign_up_at timestamp default now()
)


-- +goose Down

drop table if exists users;
