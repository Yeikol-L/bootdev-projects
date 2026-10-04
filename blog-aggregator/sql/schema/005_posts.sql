-- id - a unique identifier for the post
-- created_at - the time the record was created
-- updated_at - the time the record was last updated
-- title - the title of the post
-- url - the URL of the post (this should be unique)
-- description - the description of the post
-- published_at - the time the post was published
-- feed_id - the ID of the feed that the post cam
-- +goose up
CREATE TABLE posts (
    id  UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose down
DROP TABLE posts;
