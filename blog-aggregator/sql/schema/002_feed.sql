-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP UNIQUE NOT NULL,
    updated_at TIMESTAMP UNIQUE NOT NULL

);
-- +goose Down
DROP TABLE feeds;
