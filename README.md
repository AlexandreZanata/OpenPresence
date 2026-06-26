# OpenPresence

Biometric time and attendance platform — facial recognition, liveness detection, geofencing, and fraud detection.

**Stack:** Go + Rust backend · Kotlin Multiplatform mobile · DDD + TDD

**For coding agents:** start with **[AGENTS.md](AGENTS.md)**.

**Language policy:** 100% English — code, docs, comments, commits, and agent output.

## Documentation

Full documentation index: **[docs/README.md](docs/README.md)**

| Topic | Document |
|-------|----------|
| Product | [docs/PRODUCT-OVERVIEW.md](docs/PRODUCT-OVERVIEW.md) |
| Architecture | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) |
| Domain model | [docs/DOMAIN-MODEL.md](docs/DOMAIN-MODEL.md) |
| Business rules | [docs/BUSINESS-RULES.md](docs/BUSINESS-RULES.md) |
| API | [docs/API-CONTRACT.md](docs/API-CONTRACT.md) |
| Commits | [docs/COMMIT-CONVENTIONS.md](docs/COMMIT-CONVENTIONS.md) |
| Agent tasks | [docs/AGENT-IMPLEMENTATION-GUIDE.md](docs/AGENT-IMPLEMENTATION-GUIDE.md) |

## Quick start (developers)

```bash
pip install -r agent-harness/requirements.txt
./scripts/setup-local-ai.sh
./scripts/verify-scaffold.sh          # layout + go build/test/vet
./agent-harness/resolve-rules.sh api endpoint auth
```

Requires **Go 1.22+** for `services/attendance`.

## Project layout

```
docs/                   # Product and technical documentation
agent-rules/            # Agent Harness rule library
agent-harness/          # Rule resolution tooling
go.work                 # Go workspace (monorepo)
services/
  attendance/           # Attendance bounded context (Go)
infra/
  docker-compose.yml    # Local stack placeholder (commented)
  k8s/                  # Helm charts (planned)
  terraform/            # IaC (planned)
mobile/                 # (planned) KMP app
scripts/
  verify-scaffold.sh    # Manual layout + toolchain verification
.local/                 # Local AI tasks (gitignored)
.cursor/                # Cursor rules (gitignored)
```

## Before implementation

1. Review [docs/NEW-PROJECT-CHECKLIST.md](docs/NEW-PROJECT-CHECKLIST.md)
2. Follow [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) and [docs/COMMIT-CONVENTIONS.md](docs/COMMIT-CONVENTIONS.md)
3. Implement per [docs/AGENT-IMPLEMENTATION-GUIDE.md](docs/AGENT-IMPLEMENTATION-GUIDE.md)

## License

MIT — see [LICENSE](LICENSE). Biometric component licenses: [docs/BIOMETRICS.md](docs/BIOMETRICS.md).
