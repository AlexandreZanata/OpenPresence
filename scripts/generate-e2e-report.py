#!/usr/bin/env python3
"""Generate .local/reports/e2e-last-run.md from phase results and business-rules-coverage.json."""

from __future__ import annotations

import json
import sys
from datetime import datetime, timezone
from pathlib import Path


def main() -> int:
    if len(sys.argv) != 4:
        print("usage: generate-e2e-report.py <coverage.json> <phase-results.json> <out.md>", file=sys.stderr)
        return 2

    coverage_path = Path(sys.argv[1])
    results_path = Path(sys.argv[2])
    out_path = Path(sys.argv[3])

    coverage = json.loads(coverage_path.read_text())
    payload = json.loads(results_path.read_text())
    phase_results: dict[str, str] = payload.get("phases", {})
    overall = payload.get("overall", "UNKNOWN")
    mode = payload.get("mode", "")
    flags = payload.get("flags", "")

    def phase_ok(phases: list[str]) -> str:
        if not phases:
            return "UNKNOWN"
        states = [phase_results.get(p, "SKIP") for p in phases]
        if all(s == "PASS" for s in states):
            return "PASS"
        if any(s == "FAIL" for s in states):
            return "FAIL"
        if any(s == "PASS" for s in states):
            return "PARTIAL"
        return "SKIP"

    lines: list[str] = [
        "# E2E last run",
        "",
        f"- **Generated:** {datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}",
        f"- **Mode:** `{mode}`",
        f"- **Flags:** `{flags or 'default (continue-on-error)'}`",
        f"- **Overall:** **{overall}**",
        "",
        "## Phase results",
        "",
        "| Phase | Status |",
        "|-------|--------|",
    ]

    for phase in sorted(phase_results.keys()):
        lines.append(f"| {phase} | {phase_results[phase]} |")

    lines.extend(["", "## Business rules (28/28)", "", "| Rule | Title | Status | Phases | Notes |", "|------|-------|--------|--------|-------|"])

    automated = 0
    na_count = 0
    pass_count = 0
    fail_count = 0

    for rule in coverage["rules"]:
        rid = rule["id"]
        title = rule["title"]
        status = rule.get("status", "automated")
        phases = rule.get("phases", [])
        phase_str = ", ".join(phases)

        if status == "na":
            na_count += 1
            note = rule.get("reason", "N/A")
            lines.append(f"| {rid} | {title} | **N/A** | {phase_str} | {note} |")
            continue

        automated += 1
        result = phase_ok(phases)
        if result == "PASS":
            pass_count += 1
            lines.append(f"| {rid} | {title} | PASS | {phase_str} | |")
        elif result == "FAIL":
            fail_count += 1
            lines.append(f"| {rid} | {title} | **FAIL** | {phase_str} | phase failed |")
        else:
            lines.append(f"| {rid} | {title} | {result} | {phase_str} | |")

    lines.extend([
        "",
        "## Summary",
        "",
        f"- Automated rules: {automated} ({pass_count} pass, {fail_count} fail)",
        f"- Explicit N/A: {na_count}",
        f"- Total tracked: {len(coverage['rules'])}",
        "",
        "## Optional (not counted in 28/28)",
        "",
        "| Rule | Title | Note |",
        "|------|-------|------|",
    ])

    for opt in coverage.get("optional", []):
        lines.append(f"| {opt['id']} | {opt['title']} | {opt.get('note', '')} |")

    lines.append("")
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text("\n".join(lines) + "\n")
    print(f"Wrote {out_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
