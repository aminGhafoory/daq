-- name: CreateSession :one
INSERT INTO sessions (user_id,token_hash,created_at)
		VALUES($1,$2,$3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash=$2
RETURNING user_id;