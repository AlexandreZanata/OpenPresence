# Implementation Roadmap — Enterprise Backend

Phases **05–14** focus exclusively on **backend business logic**: organizational hierarchy (public + private), placements, schedules, punch validation, fraud, authorization, biometric model acquisition, ONNX inference, and the SubmitPunch use case.

Mobile/UI and HTTP APIs are out of scope until these domain layers are solid.

> **Implementation phases live only in `.local/phases/`** (gitignored). They are never committed to the repository.

## Phase map

| Task | Local folder | Focus | Sector coverage |
|------|----------------|--------|-----------------|
| 05 | `.local/phases/task-05-org-node-domain/` | `OrgNode` tree, invariants | Secretariats, hospitals, departments, teams |
| 06 | `.local/phases/task-06-attendance-policy/` | Policy inheritance | Public strict / private flexible presets |
| 07 | `.local/phases/task-07-employee-placement/` | Lotação | Transfer between secretariats or branches |
| 08 | `.local/phases/task-08-work-schedule/` | Time accounting BR-030–034 | 12×36 health, 8h office, split shifts |
| 09 | `.local/phases/task-09-punch-record-domain/` | Punch validation BR-010–015 | Core enterprise punch rules |
| 10 | `.local/phases/task-10-fraud-detection/` | Fraud + lockout BR-012–013 | GPS, clock, device, biometric |
| 11 | `.local/phases/task-11-hierarchy-authorization/` | ABAC subtree | Manager / HR / auditor scopes |
| 12 | `.local/phases/task-12-model-download/` | Model download script | RetinaFace, MiniFASNet, AuraFace |
| 13 | `.local/phases/task-13-onnx-inference/` | ONNX inference (Rust) | Real image processing pipeline |
| 14 | `.local/phases/task-14-punch-submission-usecase/` | SubmitPunch use case | End-to-end backend orchestration |

Tasks **00–04** (foundation) are also under `.local/phases/` when present on your machine.

## Local workspace

```bash
./scripts/setup-local-ai.sh   # .cursor/rules + .local skeleton + phases template
```

Open `.local/phases/README.md` for the active task index. Copy `_template/` to create new local phases.

## Related docs

- [ORGANIZATION.md](ORGANIZATION.md) — hierarchy examples
- [BUSINESS-RULES.md](BUSINESS-RULES.md) — BR-xxx references
- [AGENT-IMPLEMENTATION-GUIDE.md](AGENT-IMPLEMENTATION-GUIDE.md) — agent task summaries
