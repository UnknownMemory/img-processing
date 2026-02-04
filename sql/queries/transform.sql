-- name: CreateTransform :one
INSERT INTO transform (original_image, user_id)
VALUES ($1, $2)
RETURNING uuid, status;