DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'attendance_app') THEN
        CREATE ROLE attendance_app LOGIN PASSWORD 'attendance_app' NOSUPERUSER NOBYPASSRLS;
    END IF;
END
$$;

GRANT USAGE ON SCHEMA public TO attendance_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON tenants, employees, punch_records, face_embeddings TO attendance_app;
