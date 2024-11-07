-- +goose Up
CREATE TABLE
    users (
        user_id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        email varchar(255) UNIQUE NOT NULL,
        password_hash VARCHAR(72) NOT NULL
    );


-- +goose Down
DROP TABLE IF EXISTS users;