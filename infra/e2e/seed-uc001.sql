-- Fixed UUIDs for UC-001 curl/httpie smoke (Bearer e2e.<tenant>.<employee>).
INSERT INTO tenants (id, slug) VALUES
  ('11111111-1111-1111-1111-111111111111', 'uc001-e2e')
ON CONFLICT (id) DO NOTHING;

INSERT INTO employees (id, tenant_id, registration, status) VALUES
  ('22222222-2222-2222-2222-222222222222',
   '11111111-1111-1111-1111-111111111111',
   'EMP-UC001-E2E', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;
