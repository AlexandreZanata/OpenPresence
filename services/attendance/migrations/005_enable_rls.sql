ALTER TABLE employees ENABLE ROW LEVEL SECURITY;
ALTER TABLE employees FORCE ROW LEVEL SECURITY;

ALTER TABLE punch_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE punch_records FORCE ROW LEVEL SECURITY;

ALTER TABLE face_embeddings ENABLE ROW LEVEL SECURITY;
ALTER TABLE face_embeddings FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON employees
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation ON punch_records
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation ON face_embeddings
    USING (tenant_id = current_setting('app.tenant_id')::uuid);
