-- name: DeleteUser :exec
update "user"
set deleted_at = $2
where ID = $1;

-- name: CreateUser :one
insert into "user" (username,password,created_at)
values ($1,$2,$3)
returning id;

-- name: UserExists :one
SELECT * FROM "user" 
WHERE username = $1;

