-- name: CreateChirp :one
INSERT INTO chirps (body, user_id, created_at, updated_at) VALUES (
    $1,
    $2,
    now(),
    now()
) RETURNING *;
-- name: GetAllChirps :many
SELECT * FROM chirps WHERE
(sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
ORDER BY
CASE WHEN sqlc.arg('sort')::text = 'asc' THEN created_at END ASC,
CASE WHEN sqlc.arg('sort')::text = 'desc' THEN created_at END DESC;
-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1;
-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;
