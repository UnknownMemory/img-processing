-- name: CreateUser :exec
INSERT INTO users (username, password)
VALUES ($1, $2);

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;