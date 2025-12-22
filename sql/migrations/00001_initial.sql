-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username text NOT NULL,
    password text NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE images (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL REFERENCES users(id),
    url text NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;
DROP TABLE images;
