# Task 04 — Row-level security (PostgreSQL)

**Status:** done  
**Phase ID:** task-04-row-level-security

## Goal

Enable PostgreSQL RLS on tenant-scoped tables and implement `SET LOCAL app.tenant_id` pattern in Go attendance infrastructure. AGENT guide Task 04.

## Scope

**In scope:**

- sqlx migrations for RLS policies
- Go helper to run queries inside tenant-scoped transactions
- Integration test proving cross-tenant isolation

**Out of scope:**

- Full schema (minimal tables for test sufficient)
- Production Helm deployment

## Acceptance

- Integration test: tenant A cannot read tenant B rows
- All [tasks.md](tasks.md) steps `[x]`

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
