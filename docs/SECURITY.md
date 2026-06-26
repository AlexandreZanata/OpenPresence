# Security

Maps to `agent-rules/03-security/`. Biometric and LGPD constraints are project-specific additions.

## Authentication

- Passwords: **Argon2id** (m=65536, t=3, p=4 minimum)
- JWT access: **15 minutes**; refresh: **7 days**
- Device registration with revocable tokens
- MFA: required for `ORG_ADMIN` and above (v1 target)

## Authorization

RBAC + ABAC per [ORGANIZATION.md](ORGANIZATION.md). Enforced in Application layer every request.

## Biometric data (LGPD Art. 11)

- Self-hosted processing only
- Embeddings: never in REST responses; expose SHA-256 hash only
- Soft-delete on deactivation; retention policy configurable per tenant
- Audit all enrollment and access to biometric logs

## Transport

- HTTPS/mTLS client ↔ gateway
- mTLS between all internal microservices

## Multi-tenancy

- `tenant_id` on every query
- PostgreSQL RLS with transaction-scoped `SET LOCAL`
- No endpoint without tenant validation middleware

## Forbidden patterns

| Pattern | Reason |
|---------|--------|
| Business logic in HTTP handlers | Bypasses domain rules |
| ORM hiding SQL | Use sqlx (Go), explicit queries |
| Plaintext embedding storage in APIs | LGPD exposure |
| `device_time` as official punch time | Fraud vector |
| Endpoints without tenant middleware | Cross-tenant leak |

## OWASP mapping (summary)

| Risk | Mitigation |
|------|------------|
| A01 Broken Access Control | RBAC/ABAC + RLS |
| A02 Cryptographic Failures | TLS, Argon2id, AES-256-GCM local |
| A03 Injection | Parameterized sqlx, input validation at boundary |
| A04 Insecure Design | Fraud layers, audit trail |
| A05 Misconfiguration | Helm values, secrets via env |
| A07 Auth failures | Short JWT, rate limiting |
| A09 Logging failures | Immutable audit_log hypertable |

Full index: `agent-rules/03-security/README.md`

## Agentic AI (ASI01–ASI10)

AI agents operate in Application layer only: read-only by default, human confirm to persist, log all LLM calls. See `agent-rules/03-security/OWASP-AGENTIC-2026.md`.
