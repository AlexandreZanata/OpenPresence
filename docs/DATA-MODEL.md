# Data Model

PostgreSQL 16 + pgvector + TimescaleDB. Migrations via sqlx (Go). Full DDL evolves in `services/*/migrations/`.

## Core tables

### tenants

Multi-tenant root. `slug` unique, `settings` JSONB.

### org_nodes

Self-referential tree (`parent_id`). `type`, `policy` JSONB with inheritance merge.

### geofence_zones

`CIRCLE` (center + radius) or `POLYGON` (JSONB coordinates). `valid_from` / `valid_until` for temporary sites.

### employees

`tenant_id`, `org_node_id`, `registration` unique per tenant, `work_schedule` JSONB, `status`.

### face_embeddings

`vector(512)` via pgvector. IVFFlat index (`lists=100`). `active` flag for soft retention (BR-006).

### punch_records (TimescaleDB hypertable)

Partitioned on `punched_at`. Server time is official. `fraud_flags` JSONB, `sync_status` for offline.

### audit_log (TimescaleDB hypertable)

Immutable: `actor_id`, `action`, `entity_type`, `entity_id`, `old_value`, `new_value`, `ip_address`.

## Row-level security

```sql
ALTER TABLE employees ENABLE ROW LEVEL SECURITY;
ALTER TABLE punch_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE face_embeddings ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON employees
  USING (tenant_id = current_setting('app.tenant_id')::UUID);
```

Go pattern (transaction-scoped):

```go
postgres.WithTenant(ctx, db, tenantID, func(tx *sqlx.Tx) error {
    return tx.GetContext(ctx, &emp, `SELECT ... FROM employees WHERE id = $1`, id)
})
```

Uses `set_config('app.tenant_id', $1, true)` (equivalent to `SET LOCAL`). Application connects as non-superuser role `attendance_app` (see migration `006_create_app_role.sql`).

## Indexing strategy

| Table | Index |
|-------|-------|
| `face_embeddings` | IVFFlat on `embedding` (cosine) |
| `punch_records` | `(tenant_id, employee_id, punched_at DESC)` |
| `org_nodes` | `(tenant_id)`, `(parent_id)` |

## Reference DDL

See attached spec section 11 for initial `CREATE TABLE` statements. Implementation copies to versioned migrations — do not run raw DDL in production.

Task reference: [AGENT-IMPLEMENTATION-GUIDE.md](AGENT-IMPLEMENTATION-GUIDE.md) Task 04 (RLS).
