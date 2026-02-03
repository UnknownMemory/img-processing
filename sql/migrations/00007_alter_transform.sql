-- +goose Up
ALTER TABLE transform
DROP COLUMN filename,
DROP COLUMN mime,
ADD filename TEXT,
ADD mime TEXT;

-- +goose Down
ALTER TABLE transform
DROP COLUMN filename,
DROP COLUMN mime,
ADD filename TEXT NOT NULL,
ADD mime TEXT NOT NULL;