# Contributing

Thank you for contributing to OpenPresence.

## Prerequisites

1. Read [AGENTS.md](../AGENTS.md) (coding agents) or this guide (humans)
2. Complete awareness of [NEW-PROJECT-CHECKLIST.md](NEW-PROJECT-CHECKLIST.md) for your area
3. Use terms from [GLOSSARY.md](GLOSSARY.md) only

## Setup

```bash
pip install -r agent-harness/requirements.txt
./scripts/setup-local-ai.sh          # restores .cursor/ and .local/ skeleton
```

Install stack tools per [STACK.md](STACK.md).

## Workflow

1. Pick or define a use case in `docs/use-cases/`
2. Resolve agent rules: `./agent-harness/resolve-rules.sh <keywords>`
3. TDD: write failing tests first ([TESTING.md](TESTING.md))
4. Implement smallest change ([COMMIT-CONVENTIONS.md](COMMIT-CONVENTIONS.md))
5. Update glossary if new domain term introduced

## Architecture decisions

Non-trivial choices require an ADR in `docs/adr/` before implementation. Use [adr-template](../agent-rules/11-documentation-and-glossary/adr-template.md).

## Code standards

- Layered DDD: Domain never imports Infrastructure
- Size caps enforced on every change
- English only in code, comments, commits
- No business logic in HTTP handlers

## Security

Never commit secrets. Follow [SECURITY.md](SECURITY.md). Report vulnerabilities privately to maintainers.

## Documentation

When changing behavior, update:

- Business rules (`BUSINESS-RULES.md`) if rule logic changes
- API contract (`API-CONTRACT.md`) if endpoints change
- Glossary if new ubiquitous language term

## License

Contributions are MIT licensed — see [LICENSE](../LICENSE).
