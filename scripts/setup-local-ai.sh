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
mkdir -p "$ROOT/.local/tasks" "$ROOT/.local/overrides" "$ROOT/.local/phases"

ensure_phases_template() {
  local tpl="$ROOT/.local/phases/_template"
  mkdir -p "$tpl"
  if [[ ! -f "$tpl/README.md" ]]; then
    cat > "$tpl/README.md" <<'EOF'
# Task: _(short title)_

**Status:** pending | active | done  
**Phase ID:** task-NN-_(slug)_

## Goal

_(one paragraph — what “done” looks like)_

## Scope

**In scope:**

- _(bullet)_

**Out of scope:**

- _(bullet)_

## Acceptance

- All steps in [tasks.md](tasks.md) marked `[x]`
- _(measurable outcomes)_

## Agent entry

1. [official_source.md](official_source.md)
2. [tasks.md](tasks.md)
EOF
    echo "CREATED $tpl/README.md"
  fi
  if [[ ! -f "$tpl/official_source.md" ]]; then
    cat > "$tpl/official_source.md" <<'EOF'
# Official sources — _(task title)_

## Repository documentation

| Document | Path |
|----------|------|
| _(name)_ | `docs/...` |

## Agent rules

```bash
./agent-harness/resolve-rules.sh <keywords>
```

## Business rules

| ID | File |
|----|------|
| BR-xxx | `docs/BUSINESS-RULES.md` |

## Glossary terms

- _(term from docs/GLOSSARY.md)_
EOF
    echo "CREATED $tpl/official_source.md"
  fi
  if [[ ! -f "$tpl/tasks.md" ]]; then
    cat > "$tpl/tasks.md" <<'EOF'
# Tasks — _(short title)_

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh <keywords>`

## Implementation

- [ ] _(step 1)_

## Validation

- [ ] Tests pass: `_(command)_`

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
EOF
    echo "CREATED $tpl/tasks.md"
  fi
}

ensure_phases_readme() {
  local readme="$ROOT/.local/phases/README.md"
  if [[ ! -f "$readme" ]]; then
    cat > "$readme" <<'EOF'
# Implementation phases (local)

> **Gitignored** — phases never go to the repository. Each machine maintains its own copy.

## Active task

**Current:** _(set active task folder)_

## Workflow

1. Open the active task folder under `.local/phases/`.
2. Read `README.md` → `official_source.md` → `tasks.md`.
3. Run `./agent-harness/resolve-rules.sh` with keywords from `official_source.md`.
4. Mark `[x]` in `tasks.md` only after real validation.
5. Update this file when moving to the next task.

## New task

Copy `_template/` to `task-NN-short-name/` and fill the three files.

See `docs/IMPLEMENTATION-ROADMAP.md` for the high-level backend phase map.
EOF
    echo "CREATED $readme"
  fi
}

ensure_phases_template
ensure_phases_readme

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
| `phases/` | Implementation tasks (README + official_source + tasks.md) — **never in git** |
| `tasks/` | Pointer to active phase (`current.md`) |
| `context.md` | Session context, decisions, scratch notes |
| `overrides/` | Optional rules that layer on `agent-rules/` |
| `ROADMAP.md` | Implementation progress tracker |

Official rules: `agent-rules/` — do not duplicate the full tree here.

Run `./scripts/setup-local-ai.sh` after clone to restore `.cursor/rules/` and local skeleton.
EOF
        ;;
      ROADMAP.md)
        cat > "$dest" <<'EOF'
# OpenPresence — Implementation Roadmap

> Local progress tracker. One step at a time; validate before moving on.

See `docs/IMPLEMENTATION-ROADMAP.md` for the enterprise backend phase map (tasks 05–14).

Phases with checklists live in `.local/phases/` only.
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

- Active phase: `.local/phases/README.md`
- Rules: `./agent-harness/resolve-rules.sh <keywords>`
- Glossary: `docs/GLOSSARY.md`

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
  .local/         — tasks, phases, context, overrides (gitignored)

Next: open .local/phases/README.md and start the active task.

EOF
