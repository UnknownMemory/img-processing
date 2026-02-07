-- name: CreateTransform :one
INSERT INTO transform (original_image, user_id)
VALUES ($1, $2)
RETURNING uuid, status;

-- name: UpdateTransform :exec
UPDATE transform
SET status = $1, mime = $2
WHERE uuid = $3 and user_id =$4;