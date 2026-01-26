-- name: CreateImage :one
INSERT INTO images (user_id, filename, file_size, mime)
VALUES ($1, $2, $3, $4)
RETURNING uid;