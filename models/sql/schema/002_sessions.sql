-- +goose Up
CREATE TABLE
    sessions (
        user_id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        token_hash VARCHAR(255) NOT NULL UNIQUE,
        CONSTRAINT fk_sessions FOREIGN KEY (user_id) REFERENCES users (user_id)
    );

-- +goose Down
DROP TABLE IF EXISTS sessions;