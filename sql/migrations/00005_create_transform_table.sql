-- +goose Up
CREATE TABLE transform (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    original_image BIGSERIAL REFERENCES images(id),
    user_id BIGSERIAL REFERENCES users(id),
    filename TEXT NOT NULL,
    mime TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE transform;
