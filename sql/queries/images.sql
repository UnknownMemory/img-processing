-- name: CreateImage :one
INSERT INTO images (user_id, filename, file_size, mime)
VALUES ($1, $2, $3, $4)
RETURNING uid;

-- name: GetImage :one
SELECT uid, images.filename, mime, images.file_size, images.created_at
FROM images
WHERE user_id = $1 AND uid = $2;

-- name: ImageExists :one
SELECT
EXISTS(
    SELECT images.uid
    FROM images
    WHERE user_id = $1 AND uid = $2
);