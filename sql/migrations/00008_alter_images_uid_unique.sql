-- +goose Up
ALTER TABLE images
DROP uid,
ADD uid UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE;

-- +goose Down
ALTER TABLE images
DROP uid,
ADD uid UUID NOT NULL DEFAULT gen_random_uuid();