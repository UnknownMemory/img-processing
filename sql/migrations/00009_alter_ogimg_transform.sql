-- +goose Up
ALTER TABLE transform
DROP COLUMN original_image,
ADD original_image UUID REFERENCES images("uid");

-- +goose Down
ALTER TABLE transform
DROP COLUMN original_image,
ADD original_image BIGSERIAL REFERENCES images(id);