-- name: CreateImage :one
INSERT INTO images (user_id, filename, file_size, mime)
VALUES ($1, $2, $3, $4)
RETURNING uid;

-- name: GetImage :one
SELECT images.uid, images.filename, images.mime, NULL as status, images.created_at
FROM images
WHERE images.user_id = $1 AND images.uid = $2

UNION

SELECT transform.uuid, transform.filename, transform.mime, transform.status, transform.created_at
FROM transform
WHERE transform.user_id = $1 AND transform.uuid = $2;

-- name: ImageExists :one
SELECT
EXISTS(
    SELECT images.uid
    FROM images
    WHERE user_id = $1 AND uid = $2
);