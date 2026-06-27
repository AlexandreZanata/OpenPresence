-- Dev admin panel test users (passwords are mock-only until auth service exists).
-- Applied by: ./scripts/seed-dev-users.sh

INSERT INTO tenants (id, slug) VALUES
  ('11111111-1111-1111-1111-111111111111', 'dev-admin')
ON CONFLICT (id) DO NOTHING;

INSERT INTO employees (id, tenant_id, registration, status) VALUES
  ('33333333-3333-3333-3333-333333333333',
   '11111111-1111-1111-1111-111111111111', 'admin', 'ACTIVE'),
  ('34444444-4444-4444-4444-444444444444',
   '11111111-1111-1111-1111-111111111111', 'manager', 'ACTIVE'),
  ('35555555-5555-5555-5555-555555555555',
   '11111111-1111-1111-1111-111111111111', 'hr', 'ACTIVE'),
  ('36666666-6666-6666-6666-666666666666',
   '11111111-1111-1111-1111-111111111111', 'auditor', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;
