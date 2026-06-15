# v0.0.28 Phase 1 - Manager Image Packaging

## Work

- [x] Added `deploy/containerfiles/manager.Containerfile`.
- [x] Built `/manager` from `cmd/manager`.
- [x] Kept the runtime image non-root.

## Verification

- `docker build -f deploy/containerfiles/manager.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate:v0.0.28 .`

## Remaining Scope

- Confirm the image publish workflow after merge.

## Linked Issue

- #135
