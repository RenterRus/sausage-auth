-- name: Register :exec
insert into users (user_login, uuid, otp_hash, created_at, last_sign_up_at) values (sqlc.narg('login'), gen_random_uuid(), sqlc.narg('hash'), now(), now());

-- name: Confirmed :exec
update users set confirmed = true where user_login = sqlc.narg('login');

-- name: IsExist :one
select exists(select 1 from users WHERE user_login = sqlc.narg('login'));

-- name: OtpHash :one
select otp_hash from users where user_login = sqlc.narg('login');


-- name: LoginByUUID :one
select user_login from users where uuid = sqlc.narg('uuid');