CREATE TABLE "schema_migrations" (version varchar(128) primary key);
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
CREATE TABLE sessions (
    id integer PRIMARY KEY AUTOINCREMENT,
    session_token text NOT NULL UNIQUE,
    user_id integer NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    is_active integer NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20260616004011'),
  ('20260616004933');
