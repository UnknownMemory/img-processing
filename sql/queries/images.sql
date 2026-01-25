-- name: CreateImage :exec
INSERT INTO images (user_id, filename, file_size, mime, url)
VALUES ($1, $2, $3, $4, $5);