-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT feeds.*, users.name as user from feeds INNER JOIN users ON (feeds.user_id = users.id);

-- name: GetFeedByUrl :one
SELECT * FROM feeds where url = $1;

-- name: CreateFeedFollow :one
WITH feed_follows_record as (
    INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
    VALUES($1, $2, $3, $4, $5) RETURNING *
    )
select r.*, users.name as user_name, feeds.name as feed_name
from feed_follows_record as r
    INNER JOIN feeds ON (feeds.id = r.feed_id)
    INNER JOIN users ON (users.id = r.user_id);

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name as feed_name, users.name as user_name FROM feed_follows
INNER JOIN feeds ON (feeds.id = feed_follows.feed_id)
INNER JOIN users ON (users.id = feed_follows.user_id)
where feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows where user_id = $1 and feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds set last_fetched_at = now(), updated_at = now() where id = $1;

-- name: GetNextFeedToFetch :one
SELECT * from feeds ORDER BY last_fetched_at ASC NULLS FIRST limit 1;
