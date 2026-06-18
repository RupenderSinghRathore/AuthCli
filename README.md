# AuthCLI

AuthCLI is a secure command-line login system written in Go. It supports user registration, username/password authentication, optional TOTP-based 2FA, session persistence, SQLite storage, and dbmate-managed schema migrations. The project also includes a Docker-based workflow for running the app in a container.

## Features

- Interactive prompt with history and tab-completion.
- User registration with unique usernames
- Password hashing with bcrypt
- Login with username and password
- Optional TOTP-based 2FA
- Session creation and session resume when the session file is preserved
- SQLite persistence
- Schema migrations managed with dbmate
- Docker and Docker Compose support

## Tech Stack

- Go
- SQLite
- dbmate
- sqlc
- Docker

## Project Structure

```text
.
├── cmd/app/             # CLI entrypoint and command handlers
├── db/migrations/       # dbmate migration files
├── db/queries/          # SQL queries used by sqlc
├── internal/database/   # sqlc-generated database layer
├── Dockerfile
├── compose.yaml
└── docker-entrypoint.sh
```

## Database Design

```text
┌──────────────────────┐
│        users         │
├──────────────────────┤
│ id PK                │
│ username UNIQUE      │
│ password_hash        │
│ mfa_enabled          │
│ totp_secret          │
│ failed_attempts      │
│ locked_until         │
│ created_at           │
│ last_login_at        │
└──────────┬───────────┘
           │1
           │
           │
           │1
┌──────────▼───────────┐
│      sessions        │
├──────────────────────┤
│ id PK                │
│ session_token        │
│ user_id FK           │
│ created_at           │
│ expires_at           │
│ is_active            │
└──────────────────────┘
```

## Environment Variables

The application uses these environment variables:

- `DATABASE_PATH`: SQLite file path used by the Go application
- `DATABASE_URL`: SQLite connection URL used by dbmate

Example local `.env`:

```env
DATABASE_URL=sqlite:///home/your-user/path/to/AuthCli/db/myapp.db
DATABASE_PATH=/home/your-user/path/to/AuthCli/db/myapp.db
```

## Run Locally

### 1. Install dependencies

- Go
- SQLite
- dbmate

### 2. Run migrations

```bash
dbmate up
```

### 3. Start the CLI

```bash
go run ./cmd/app
```

## Run With Docker

The Docker setup builds the Go binary, installs dbmate, runs migrations on container startup, and then launches the interactive CLI.

### Build the image

```bash
docker build -t authcli .
```

### Run the container directly

```bash
docker run --rm -it \
  -e DATABASE_PATH=/app/db/myapp.db \
  -e DATABASE_URL=sqlite:///app/db/myapp.db \
  -v "$(pwd)/db:/app/db" \
  authcli
```

## One-Command Docker Compose Flow

For this project, `docker compose run` is the better fit than `docker compose up` because the application is an interactive CLI, not a long-running background service.

### Run with Docker Compose

```bash
docker compose run --build --rm authcli
```

If your machine uses the older Compose binary:

```bash
docker-compose run --build --rm authcli
```

### Persistence behavior

The Compose setup mounts the local `db/` directory into the container:

```yaml
volumes:
  - ./db:/app/db
```

That means:

- database changes persist across container runs
- if `db/myapp.db` already exists on your machine, the container uses that same file
- if `db/myapp.db` does not exist yet, dbmate creates it and applies migrations

dbmate migrates the database schema forward. It does not erase existing records unless you explicitly remove or replace the database file.

## Session Behavior

After a successful login or registration, the app creates a session token and stores it under:

```text
~/.config/AuthCLI/session
```

On the next launch, the app attempts to restore the logged-in user from that session.

For local runs, this file lives on your machine and is reused automatically.

For Docker runs, both the database and the session file persist through bind mounts. The Compose setup mounts `./db` to `/app/db` for SQLite data and `./.docker-config` to `/root/.config` so the CLI session file survives across `docker compose run --rm` executions.

## Available Commands

When logged out:

```text
Available Commands

Authentication:
  register      Create a new user account
  login         Login with username and password
  help          Show this help message
  exit          Quit the program
```

When logged in:

```text
Available Commands

Account:
  whoami        Show current user details
  enable-2fa    Enable TOTP-based MFA
  disable-2fa   Disable MFA
  logout        End current session
  help          Show this help message
  exit          Quit the program
```

## Authentication Flow

### Register

1. Create a username and password
2. The password is hashed with bcrypt before storage
3. A session is created immediately after registration

### Login

1. Enter username and password
2. If 2FA is enabled, enter the current TOTP code
3. A new session is created on successful authentication

### 2FA Setup

1. Log in
2. Run `enable-2fa`
3. Copy the displayed secret into your authenticator app
4. Enter the current TOTP code to confirm setup

## Security Notes

- Passwords are stored as bcrypt hashes, not plaintext
- Usernames must be 3 to 32 characters and may contain letters, numbers, `_`, and `-`
- Passwords are limited to bcrypt-compatible lengths
- Repeated failed logins increase the failed attempt counter
- Accounts are temporarily locked after 3 failed login attempts
- Sessions are stored separately from passwords and tied to the database

## Example Workflow

```text
> register
username: alice
Password:
registered successfully!

> whoami
username: alice
2FA: disabled

> enable-2fa
use your authenticator app to setup totp with these values:
code name: alice
your key: XXXXXXXX
your authenticator key: 123456
enabled TOTP based 2FA

> logout
logged out successfully
```
