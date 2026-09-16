-- name: GetRefreshToken :one
select r.refresh_hash, (expired_at <= now()) as is_expired, r.user_agent, r.user_login,
exists(select 1 from blacklist_refresh b where b.refresh_hash = r.refresh_hash) as block 
from refreshlist r where r.user_login = sqlc.narg('user_login') and r.refresh_hash = sqlc.narg('hash') and r.user_agent = sqlc.narg('user_agent');

-- name: SetBlockRefresh :exec
insert into blacklist_refresh (refresh_hash) values (sqlc.narg('refresh_hash'));

-- name: SetRefreshHash :exec
insert into refreshlist(refresh_hash, user_login, user_agent) values(sqlc.narg('refresh_hash'), sqlc.narg('user_login'), sqlc.narg('user_agent'));

-- name: RemoveRefreshByHash :exec
delete from refreshlist where refresh_hash = sqlc.narg('refresh_hash');

-- name: RemoveRefreshByLogin :many
delete from refreshlist where user_login = sqlc.narg('user_login') returning refresh_hash;

-- name: GetUUIDByLogin :one
select uuid from users WHERE user_login = sqlc.narg('login');

-- name: RemoveOldRefresh :many
delete from refreshlist where refresh_hash = sqlc.narg('refresh_hash') and user_agent = sqlc.narg('user_agent') returning refresh_hash, user_login;

-- name: RemoveOldRefreshByLoginUA :many
delete from refreshlist where refresh_hash = sqlc.narg('user_login') and user_agent = sqlc.narg('user_agent') returning refresh_hash, user_login;


-- name: UpdateLastSighUp :exec
update users set last_sign_up_at = now() where user_login = sqlc.narg('user_login');