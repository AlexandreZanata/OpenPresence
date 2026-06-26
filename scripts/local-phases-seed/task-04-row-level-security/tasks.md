# Tasks — Row-level security

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh owasp authz tenant migration`
- [ ] Docker available for Postgres testcontainer

## Migrations

- [ ] Add `services/attendance/migrations/` with sqlx
- [ ] Migration: create minimal `employees` table with `tenant_id`
- [ ] Migration: `ALTER TABLE employees ENABLE ROW LEVEL SECURITY`
- [ ] Migration: `CREATE POLICY tenant_isolation ON employees USING (...)`
- [ ] Repeat for `punch_records` and `face_embeddings` stubs per `docs/DATA-MODEL.md`

## Go infrastructure

- [ ] `internal/infrastructure/postgres/tenant_tx.go` — `WithTenant(ctx, tenantID, fn)`
- [ ] Executes `SET LOCAL app.tenant_id = $1` inside transaction
- [ ] Repository example: `GetEmployee` uses `WithTenant`

## Integration tests

- [ ] testcontainers-go Postgres 16
- [ ] Seed tenant A and tenant B employees
- [ ] Assert tenant A session cannot SELECT tenant B row
- [ ] Assert same query with correct tenant returns row

## Validation

- [ ] `go test ./...` including integration tag if used: `-tags=integration`
- [ ] Migrations apply cleanly on empty DB
- [ ] Document `DATABASE_URL` in `.env.example` if new vars needed

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
