CREATE TABLE punch_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants (id),
    employee_id UUID NOT NULL REFERENCES employees (id),
    punch_type TEXT NOT NULL,
    punched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    sync_status TEXT NOT NULL DEFAULT 'SYNCED'
);

CREATE INDEX idx_punch_records_tenant ON punch_records (tenant_id, employee_id, punched_at DESC);
