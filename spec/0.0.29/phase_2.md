# v0.0.29 Phase 2 - Bundle Image Publishing

## Work

- [x] Added `deploy/containerfiles/bundle.Containerfile`.
- [x] Added `.github/workflows/operator-bundle.yml`.
- [x] Configured publish of `v0.0.29` and `latest` bundle tags.

## Verification

- `docker build -f deploy/containerfiles/bundle.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.29 .`

## Remaining Scope

- Confirm post-merge workflow publish.

## Linked Issue

- #141
