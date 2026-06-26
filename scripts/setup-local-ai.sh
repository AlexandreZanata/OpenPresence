#!/usr/bin/env bash
# Restore gitignored local AI config after clone: .cursor/rules/ and .local/ skeleton.
#
# Usage: ./scripts/setup-local-ai.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUNDLE="$ROOT/agent-harness/cursor-rules-bundle"

# --- .cursor/rules from harness bundle ---
if [[ ! -d "$BUNDLE" ]]; then
  echo "ERROR: cursor-rules-bundle missing at $BUNDLE" >&2
  echo "Re-run harness install or copy from GoodPraticesForLLMSandAgents." >&2
  exit 1
fi

mkdir -p "$ROOT/.cursor/rules"
if command -v rsync >/dev/null 2>&1; then
  rsync -a "$BUNDLE/" "$ROOT/.cursor/rules/"
else
  cp -a "$BUNDLE/." "$ROOT/.cursor/rules/"
fi
echo "RESTORED $ROOT/.cursor/rules/"

# --- .local skeleton (never overwrite existing files) ---
mkdir -p "$ROOT/.local/tasks" "$ROOT/.local/overrides"

for f in README.md ROADMAP.md tasks/README.md tasks/current.md context.md overrides/README.md; do
  dest="$ROOT/.local/$f"
  if [[ ! -f "$dest" ]]; then
    mkdir -p "$(dirname "$dest")"
    case "$f" in
      README.md)
        cat > "$dest" <<'EOF'
# Local AI workspace (gitignored)

Machine-local files for coding agents. Not versioned — each developer maintains their own copy.

| Path | Purpose |
|------|---------|
| `tasks/` | Active and backlog tasks for AI sessions |
| `context.md` | Session context, decisions, scratch notes |
| `overrides/` | Optional rules that layer on `agent-rules/` |
| `ROADMAP.md` | Implementation progress tracker |

Official rules: `agent-rules/` — do not duplicate the full tree here.
EOF
        ;;
      ROADMAP.md)
        cat > "$dest" <<'EOF'
# OpenPresence — Implementation Roadmap

> Local progress tracker. One step at a time; validate before moving on.

## Phase 0 — Project bootstrap

- [x] Agent Harness installed (`agent-rules/`, `agent-harness/`)
- [x] Docs templates created (`docs/`)
- [x] Local AI workspace (`.local/`, `.cursor/`)
- [ ] NEW-PROJECT-CHECKLIST completed
- [ ] Glossary and API contract filled
- [ ] Stack and repo structure decided

## Phase 1 — Domain definition

- [ ] Entities, aggregates, and value objects documented
- [ ] Business rules (GIVEN/WHEN/THEN)
- [ ] State machines for multi-status entities
- [ ] Use cases in `docs/use-cases/`

## Phase 2 — Implementation

- [ ] _(add milestones as the project evolves)_

---

See `docs/NEW-PROJECT-CHECKLIST.md` for the full pre-coding checklist.
EOF
        ;;
      tasks/README.md)
        cat > "$dest" <<'EOF'
# AI task files

| File | Purpose |
|------|---------|
| `current.md` | Active task — agent reads this first for focused work |
| `backlog.md` | Optional queued tasks |

**Workflow:** Update `current.md` at session start. Clear or archive when done.
EOF
        ;;
      tasks/current.md)
        cat > "$dest" <<'EOF'
# Current task

> Update this file at the start of each AI session. Delete sections when done.

## Goal

_(what to accomplish in this session)_

## Context

- Rules: `./agent-harness/resolve-rules.sh <keywords>`
- Glossary: `docs/GLOSSARY.md`
- Checklist: `docs/NEW-PROJECT-CHECKLIST.md`

## Acceptance criteria

- [ ] _(measurable outcome)_

## Notes

_(scratch space)_
EOF
        ;;
      context.md)
        cat > "$dest" <<'EOF'
# Session context (local)

Scratch pad for decisions, open questions, and links. Agents may read this for continuity.

## Open questions

- _(list items the agent must ASK about — never assume)_

## Decisions

| Date | Decision | Rationale |
|------|----------|-----------|
| | | |
EOF
        ;;
      overrides/README.md)
        cat > "$dest" <<'EOF'
# Project-specific rule overrides

Optional `.md` rule files layer on top of `agent-rules/` — not a replacement.

Example: `overrides/team-conventions.md`

Tell the agent to load overrides explicitly; they are not in `manifest.yaml`.
EOF
        ;;
    esac
    echo "CREATED $dest"
  else
    echo "SKIP $dest (exists)"
  fi
done

cat <<EOF

Local AI setup complete.

  .cursor/rules/  — Cursor validation rules (from harness bundle)
  .local/         — tasks, context, overrides

Next: fill docs/NEW-PROJECT-CHECKLIST.md before writing application code.

EOF
