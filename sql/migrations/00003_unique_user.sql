-- +goose Up
ALTER TABLE users
ADD CONSTRAINT u_username UNIQUE (username);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT IF EXISTS u_username;
