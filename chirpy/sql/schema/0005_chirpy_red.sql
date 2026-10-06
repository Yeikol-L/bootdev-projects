-- +goose up
ALTER TABLE users ADD COLUMN is_chirpy_red bool NOT NULL DEFAULT false;
-- +goose down
ALTER TABLE users DROP COLUMN is_chirpy_red;
