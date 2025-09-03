!/usr/bin/env bash
set -euo pipefail

CONTAINER="agrobot-db-1"
DB="${1:-appdb}"
USER="app"
PGPASSWORD="secret"
DUMP_FILE="${2:-dump_$(date +%Y%m%d_%H%M%S).sql}"

docker exec -e PGPASSWORD="$PGPASSWORD" "$CONTAINER" \
    pg_dump -U "$USER" -d "$DB" -Fp --no-owner --no-privileges \
    > "$DUMP_FILE"

echo "Dumped $DB to $DUMP_FILE"

