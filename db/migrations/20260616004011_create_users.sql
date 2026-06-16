-- migrate:up

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    mfa_enabled INTEGER NOT NULL DEFAULT 0,
    totp_secret TEXT,

    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until DATETIME,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at DATETIME
);

-- migrate:down

DROP TABLE users;
