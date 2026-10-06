-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password,created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, now(), now())
RETURNING *;
-- name: DeleteAllUsers :exec
DELETE FROM users;
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;
-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3 WHERE id = $1 RETURNING *;
-- name: UpgradeUser :one
UPDATE users SET is_chirpy_red = true, updated_at = now() WHERE id = $1 RETURNING *;
