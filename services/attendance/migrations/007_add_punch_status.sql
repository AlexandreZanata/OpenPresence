ALTER TABLE punch_records
    ADD COLUMN status TEXT NOT NULL DEFAULT 'VALID';

CREATE INDEX idx_punch_records_employee_status
    ON punch_records (tenant_id, employee_id, status, punched_at DESC);
