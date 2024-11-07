-- name: CreateUser :exec
INSERT INTO 
    users(
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
        );        