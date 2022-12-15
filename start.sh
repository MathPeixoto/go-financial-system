#!/bin/sh

set -e

echo "run db migrations"
source /app/app.env
/app/migrate -path /app/migration -database "$DATABASE_SOURCE" -verbose up

echo "run server"
exec "$@"