-- name: CreateSession :one
INSERT INTO sessions (session_token, user_id, expires_at)
    VALUES (?, ?, datetime ('now', '+7 days'))
RETURNING
    *;

-- name: GetActiveSession :one
SELECT
    *
FROM
    sessions
WHERE
    session_token = ?
    AND is_active = 1
    AND expires_at > CURRENT_TIMESTAMP;

