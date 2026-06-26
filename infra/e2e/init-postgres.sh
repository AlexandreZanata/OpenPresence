#!/bin/bash
set -euo pipefail

for f in /migrations/*.sql; do
  psql -v ON_ERROR_STOP=1 -U openpresence -d openpresence -f "$f"
done
