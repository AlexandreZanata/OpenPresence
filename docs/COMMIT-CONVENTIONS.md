# Commit Conventions

Aligned with [agent-rules/00-core/change-discipline.md](../agent-rules/00-core/change-discipline.md).

## Format

[Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

| Type | Use |
|------|-----|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code change, no behavior change |
| `test` | Tests only |
| `docs` | Documentation only |
| `chore` | Tooling, deps, CI |
| `perf` | Performance improvement |

### Scopes (examples)

`attendance`, `biometric`, `organization`, `workforce`, `mobile`, `infra`, `docs`

### Subject rules

- English only
- Imperative mood: "add geofence validator" not "added"
- Max 72 characters
- No period at end

### Body

- Explain **why**, not what (the diff shows what)
- Reference business rules: `Implements BR-020 circle geofence check`
- Reference ADRs when applicable: `See ADR-002`

## One logical change per commit

Each commit MUST be exactly one of:

- Feature implementation
- Bug fix
- Refactor (no behavior change)
- Test-only change
- Documentation-only change

**Never** mix refactor + feature in the same commit.

## Edit order (recommended)

1. Tests or contract (TDD / API change)
2. Domain layer
3. Application layer
4. Infrastructure adapters
5. Interfaces (HTTP, gRPC, UI)
6. Glossary / docs if domain terms changed

## Before committing

- [ ] Scope matches request
- [ ] Relevant tests run and pass
- [ ] Glossary updated if new domain term
- [ ] No secrets, PII, or `.env` committed
- [ ] English only in strings and comments
- [ ] Size caps: ≤80 lines/function, ≤200 lines/file

## Examples

```
feat(attendance): add circle geofence validation

Implements BR-020. Haversine distance with allowedDeviation buffer.
TDD: tests written first in internal/domain/geofence.
```

```
docs: add domain model and business rules

Product spec from initial design. Fills NEW-PROJECT-CHECKLIST
architecture and domain sections.
```

```
fix(biometric): load ONNX sessions once at startup

Prevents per-request model load. Target P99 < 500ms for punch path.
```

## Pull requests

- Title follows commit convention
- One logical change per PR (prefer small PRs)
- Link related use case or BR id in description
- CI must be green before merge
