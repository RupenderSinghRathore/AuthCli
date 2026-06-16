-- name: CreateUser :one
INSERT INTO users (username, password_hash)
    VALUES (?, ?)
RETURNING
    *;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    username = ?;

