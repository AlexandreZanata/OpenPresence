# AGENTS.md — OpenPresence Entry Point for Coding Agents

> **Read this first** in any new agent session (Cursor, Claude Code, Codex, Windsurf, etc.).

**Language:** 100% English — code, comments, docs, commits, and all agent output.

---

## What this repo is

| Is | Is not |
|----|--------|
| OpenPresence application (presence / attendance platform) | The Agent Harness source repo |
| AI-assisted with enterprise agent rules | Ready-to-run without checklist completion |
| Layered DDD + OWASP-aligned security | A place to invent domain terms without glossary |

When rules conflict with existing code, **rules prevail** — unless the user explicitly overrides for a task.

---

## Rules path (resolve first)

```bash
pip install -r agent-harness/requirements.txt   # once per machine
./agent-harness/rules-path.sh                     # → agent-rules/
```

| Config | `rules_dir` |
|--------|-------------|
| `agent-harness/harness.config.yaml` | `agent-rules/` |

Never hardcode `rules/` — this project uses `agent-rules/`.

---

## Local AI workspace (gitignored)

| Path | Purpose |
|------|---------|
| `.local/tasks/current.md` | Active task — read at session start |
| `.local/context.md` | Decisions and open questions |
| `.local/overrides/` | Optional rules on top of `agent-rules/` |
| `.cursor/rules/` | Cursor validation rules |

After clone, restore local config:

```bash
./scripts/setup-local-ai.sh
```

---

## Always load (base context)

1. `agent-rules/AGENT-CORE-PRINCIPLES.md` — architecture contract
2. `agent-rules/00-core/size-and-complexity-limits.md` — **80 lines/function, 200 lines/file, cyclomatic ≤10**
3. `agent-rules/09-ai-agent-specific/token-economy.md`
4. `agent-rules/09-ai-agent-specific/anti-hallucination.md`

Cursor: `.cursor/rules/*.mdc` applies automatically (restore via `setup-local-ai.sh`).

### Ponytail (YAGNI)

```bash
./agent-harness/resolve-rules.sh yagni minimal ponytail
```

Harness security and TDD rules **override** Ponytail when they conflict.

---

## Conditional load (task-specific)

Load **2–6 files only** — not the entire rule tree.

```bash
./agent-harness/resolve-rules.sh <keywords from task>
```

| Task | Example keywords |
|------|------------------|
| API endpoint | `api endpoint auth validation contract` |
| Security review | `owasp security authz bola agentic` |
| Domain feature | `domain layer state event` |
| Bug fix | `bugfix regression error` |

### Cursor: task-scoped rule file (optional)

```bash
./agent-harness/generate-task-rules.sh api endpoint auth
./agent-harness/generate-task-rules.sh --clean
```

---

## Agent protocol

1. Run `./agent-harness/rules-path.sh` → `agent-rules/`.
2. Read `.local/tasks/current.md` if it exists.
3. Identify task keywords → run `resolve-rules.sh`.
4. State which rule files you loaded (brief list).
5. **ASK** if `docs/NEW-PROJECT-CHECKLIST.md` items are blank — never assume business rules.
6. Smallest diff; one logical change per commit.
7. Verify after each edit — do not claim tests passed without running them.

---

## Key references

| Document | Purpose |
|----------|---------|
| `docs/NEW-PROJECT-CHECKLIST.md` | Complete before first line of code |
| `docs/GLOSSARY.md` | Ubiquitous language |
| `docs/API-CONTRACT.md` | API sketch |
| `agent-rules/AGENT-CORE-PRINCIPLES.md` | Full architecture contract |
| `agent-rules/STRUCTURE.md` | Rule tree + task mapping |
| `agent-rules/03-security/README.md` | Security index |
| `THIRD_PARTY_NOTICES.md` | Ponytail MIT attribution |

---

## Optional local overrides

Project-specific rules: `.local/overrides/` (gitignored). Harness rules still apply unless user explicitly waives.
