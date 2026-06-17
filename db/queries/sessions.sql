-- name: CreateSession :one
INSERT INTO sessions (session_token, user_id, expires_at)
    VALUES (lower(hex (randomblob (32))), ?, ?)
RETURNING
    session_token;

-- name: GetActiveSession :one
SELECT
    *
FROM
    sessions
WHERE
    session_token = ?
    AND is_active = 1
    AND expires_at > CURRENT_TIMESTAMP;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= CURRENT_TIMESTAMP;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE user_id = ?;
