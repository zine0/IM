-- name: DeleteUser :exec
update "user"
set deleted_at = $2
where ID = $1;

-- name: CreateUser :exec
insert into "user" (username,password,created_at)
values ($1,$2,$3);

-- name: UserExists :one
SELECT EXISTS(
    SELECT 1 FROM "user" WHERE username = $1
) AS exists;
