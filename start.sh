#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migration -database "$DATABASE_SOURCE" -verbose up

echo "run server"
exec "$@"