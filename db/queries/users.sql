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

-- name: GetUserBySessionToken :one
SELECT
    u.*
FROM
    users u
JOIN
    sessions s ON s.user_id = u.id
WHERE
    s.session_token = ?
    AND s.is_active = 1
    AND s.expires_at > CURRENT_TIMESTAMP;
