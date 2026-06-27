#!/usr/bin/env bash
# Seed dev admin test users into local Postgres (employees + fixed UUIDs).
#
# Usage: ./scripts/seed-dev-users.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="${ROOT}/infra/docker-compose.e2e.yml"
SEED="${ROOT}/infra/dev/seed-admin-users.sql"

if [[ ! -f "$SEED" ]]; then
  echo "missing seed file: $SEED" >&2
  exit 1
fi

if ! docker compose -f "$COMPOSE_FILE" ps postgres -q 2>/dev/null | grep -q .; then
  echo "Postgres not running. Start: ./scripts/dev-backend.sh start" >&2
  exit 1
fi

docker compose -f "$COMPOSE_FILE" exec -T postgres \
  psql -U openpresence -d openpresence -v ON_ERROR_STOP=1 -f - < "$SEED"

echo "Dev admin users seeded (admin, manager, hr, auditor)."
