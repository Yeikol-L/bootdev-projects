-- +goose up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP UNIQUE NOT NULL,
    updated_at TIMESTAMP UNIQUE NOT NULL,
    CONSTRAINT feed_follows_feed_id_user_id_unique UNIQUE(feed_id, user_id)
);
-- +goose down
DROP TABLE feed_follows;
