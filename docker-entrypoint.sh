#!/bin/sh
set -eu

DATABASE_PATH="${DATABASE_PATH:-/app/db/myapp.db}"
DATABASE_URL="${DATABASE_URL:-sqlite://${DATABASE_PATH}}"

export DATABASE_PATH DATABASE_URL

mkdir -p "$(dirname "$DATABASE_PATH")"

dbmate --migrations-dir /app/db/migrations up

exec "$@"
