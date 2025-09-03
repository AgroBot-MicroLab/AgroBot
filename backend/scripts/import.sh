#!/usr/bin/env bash
set -euo pipefail

CONTAINER="agrobot-db-1"
DB="${1:-appdb}"
USER="app"
PGPASSWORD="secret"
SQL_FILE="${2:?path to .sql file required}"

[[ -r "$SQL_FILE" ]] || { echo "File not readable: $SQL_FILE"; exit 1; }

docker exec -i -e PGPASSWORD="$PGPASSWORD" "$CONTAINER" \
    psql -U "$USER" -d "$DB" -v ON_ERROR_STOP=1 < "$SQL_FILE"

echo "Imported $SQL_FILE into $DB"

