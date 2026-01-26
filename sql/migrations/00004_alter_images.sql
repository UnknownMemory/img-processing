-- +goose Up
ALTER TABLE images
DROP COLUMN url;

-- +goose Down
ALTER TABLE images
ADD url TEXT NOT NULL;
