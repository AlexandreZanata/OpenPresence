# Implementation Roadmap — Enterprise Backend

Phases **05–14** focus exclusively on **backend business logic**: organizational hierarchy (public + private), placements, schedules, punch validation, fraud, authorization, biometric model acquisition, ONNX inference, and the SubmitPunch use case.

Mobile/UI and HTTP APIs are out of scope until these domain layers are solid.

## Phase map

| Task | Focus | Sector coverage |
|------|--------|-----------------|
| [05](../scripts/local-phases-seed/task-05-org-node-domain/) | `OrgNode` tree, invariants | Secretariats, hospitals, departments, teams |
| [06](../scripts/local-phases-seed/task-06-attendance-policy/) | Policy inheritance | Public strict / private flexible presets |
| [07](../scripts/local-phases-seed/task-07-employee-placement/) | Lotação | Transfer between secretariats or branches |
| [08](../scripts/local-phases-seed/task-08-work-schedule/) | Time accounting BR-030–034 | 12×36 health, 8h office, split shifts |
| [09](../scripts/local-phases-seed/task-09-punch-record-domain/) | Punch validation BR-010–015 | Core enterprise punch rules |
| [10](../scripts/local-phases-seed/task-10-fraud-detection/) | Fraud + lockout BR-012–013 | GPS, clock, device, biometric |
| [11](../scripts/local-phases-seed/task-11-hierarchy-authorization/) | ABAC subtree | Manager / HR / auditor scopes |
| [12](../scripts/local-phases-seed/task-12-model-download/) | Model download script | RetinaFace, MiniFASNet, AuraFace |
| [13](../scripts/local-phases-seed/task-13-onnx-inference/) | ONNX inference (Rust) | Real image processing pipeline |
| [14](../scripts/local-phases-seed/task-14-punch-submission-usecase/) | SubmitPunch use case | End-to-end backend orchestration |

## Restore local phases

```bash
./scripts/setup-local-ai.sh
```

New task folders are copied to `.local/phases/` without overwriting existing files. To refresh an updated `README.md` index, merge manually or delete the local copy and re-run setup.

## Related docs

- [ORGANIZATION.md](ORGANIZATION.md) — hierarchy examples
- [BUSINESS-RULES.md](BUSINESS-RULES.md) — BR-xxx references
- [AGENT-IMPLEMENTATION-GUIDE.md](AGENT-IMPLEMENTATION-GUIDE.md) — agent task details
