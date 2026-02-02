-- +goose Up
ALTER TABLE transform
ADD status TEXT NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE transform
DROP COLUMN status;
