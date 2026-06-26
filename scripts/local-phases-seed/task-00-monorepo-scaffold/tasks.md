# Tasks — Monorepo scaffold

## Preparation

- [x] Read [README.md](README.md) and [official_source.md](official_source.md)
- [x] Run `./agent-harness/resolve-rules.sh architecture layering domain`
- [x] Confirm Go >= 1.22 installed

## Implementation

- [x] Create `services/attendance/go.mod` with module path `github.com/AlexandreZanata/OpenPresence/services/attendance`
- [x] Create directory skeleton:
  - `internal/domain/`
  - `internal/application/`
  - `internal/infrastructure/`
  - `internal/interfaces/`
- [x] Add `services/attendance/README.md` describing service responsibility
- [x] Create minimal `infra/docker-compose.yml` placeholder (commented services only)
- [x] Update root `README.md` project layout

## Validation

- [x] `cd services/attendance && go build ./...` succeeds
- [x] `cd services/attendance && go test ./...` succeeds
- [x] `./scripts/verify-scaffold.sh` — manual verification (build, test, vet, layout)
- [x] No secrets or `.env` committed

## Completion

- [x] All steps above marked `[x]`
- [x] Set next active task in `.local/phases/README.md`
