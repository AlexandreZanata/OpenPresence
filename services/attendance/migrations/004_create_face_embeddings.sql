CREATE TABLE face_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants (id),
    employee_id UUID NOT NULL REFERENCES employees (id),
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_face_embeddings_tenant ON face_embeddings (tenant_id);
