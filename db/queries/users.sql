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
    JOIN sessions s ON s.user_id = u.id
WHERE
    s.session_token = ?
    AND s.is_active = 1
    AND s.expires_at > CURRENT_TIMESTAMP;

-- name: GetUserForLogin :one
SELECT
    *
FROM
    users
WHERE
    username = ?
    AND (locked_until IS NULL
        OR locked_until <= CURRENT_TIMESTAMP);

-- name: GetUserByUsername :one
SELECT
    *
FROM
    users
WHERE
    username = ?;

-- name: RecordSuccessfulLogin :exec
UPDATE
    users
SET
    last_login_at = CURRENT_TIMESTAMP,
    failed_attempts = 0,
    locked_until = NULL
WHERE
    id = ?;

-- name: RecordFailedLogin :one
UPDATE
    users
SET
    failed_attempts = CASE WHEN failed_attempts + 1 >= sqlc.arg (max_attempts) THEN
        0
    ELSE
        failed_attempts + 1
    END,
    locked_until = CASE WHEN failed_attempts + 1 >= sqlc.arg (max_attempts) THEN
        sqlc.arg (locked_until)
    ELSE
        locked_until
    END
WHERE
    id = sqlc.arg (user_id)
RETURNING
    *;

-- name: EnableMFA :exec
UPDATE
    users
SET
    mfa_enabled = TRUE,
    totp_secret = ?
WHERE
    id = ?;

-- name: DisableMFA :exec
UPDATE
    users
SET
    mfa_enabled = 0,
    totp_secret = NULL
WHERE
    id = ?;

