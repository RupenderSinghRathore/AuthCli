-- migrate:up

CREATE TABLE sessions (
    user_id INTEGER PRIMARY KEY,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,

    is_active INTEGER NOT NULL DEFAULT 1,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- migrate:down

DROP TABLE sessions;
