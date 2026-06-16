-- migrate:up
CREATE TABLE sessions (
    id integer PRIMARY KEY AUTOINCREMENT,
    session_token text NOT NULL UNIQUE,
    user_id integer NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    is_active integer NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- migrate:down
DROP TABLE sessions;

