-- name: CreateUser :one
INSERT INTO users(
    user_id,
    created_at,
    updated_at,
    email,
    password_hash
    )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
    )
RETURNING *;
-- name: UserByEmail :one
SELECT
    user_id,
    password_hash
FROM
    users
WHERE email =$1;

-- name: UserBySession :one
SELECT users.user_id,
       email,
       password_hash
	FROM sessions
	INNER JOIN users ON
		sessions.user_id = users.user_id
	WHERE token_hash = $1;

-- name: DeleteUserSession :exec
DELETE 
    FROM sessions 
    WHERE token_hash=$1;

