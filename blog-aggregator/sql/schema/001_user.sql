-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP UNIQUE NOT NULL,
    updated_at TIMESTAMP UNIQUE NOT NULL
);
-- +goose Down
DROP TABLE users;
