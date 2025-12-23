-- +goose Up
ALTER TABLE images
DROP user_id,
ADD user_id BIGSERIAL REFERENCES users(id) ON DELETE CASCADE,
ADD uid UUID NOT NULL DEFAULT gen_random_uuid(),
ADD filename TEXT NOT NULL,
ADD mime TEXT NOT NULL,
ADD file_size BIGINT;

-- +goose Down
ALTER TABLE images
DROP user_id,
ADD user_id BIGSERIAL REFERENCES users(id),
DROP COLUMN uid,
DROP COLUMN filename,
DROP COLUMN mime,
DROP COLUMN file_size;