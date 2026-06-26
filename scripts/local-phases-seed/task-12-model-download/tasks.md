# Tasks — Biometric model download

## Preparation

- [ ] Read [README.md](README.md) and [official_source.md](official_source.md)
- [ ] Run `./agent-harness/resolve-rules.sh security biometric integrity`
- [ ] Confirm `models/` in `.gitignore`

## Manifest

- [ ] Create `models/MANIFEST.json`: file name, URL, sha256, license note
- [ ] Document expected layout in `docs/BIOMETRICS.md`

## Download script

- [ ] `scripts/download-models.sh` — curl/wget with retry
- [ ] Verify SHA-256 per file against manifest
- [ ] Idempotent: skip if hash matches
- [ ] Fail on partial download

## Verification

- [ ] `scripts/verify-models.sh` — all files present + hashes valid
- [ ] Manual run on clean `models/` directory
- [ ] Update root `README.md` quick start

## Completion

- [ ] All steps above marked `[x]`
- [ ] Update `.local/phases/README.md` active task
