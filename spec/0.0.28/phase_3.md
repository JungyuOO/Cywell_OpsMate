# v0.0.28 Phase 3 - OpenShift Manifest Image Pin

## Work

- [x] Updated `config/manager/manager.yaml` to use `ghcr.io/jungyuoo/cywell-opsmate:v0.0.28`.
- [x] Kept the existing `/manager` command and non-root security context.

## Verification

- `kubectl kustomize config/default`
- `go test ./...`

## Remaining Scope

- Apply merged manifests to OpenShift after the image publish workflow succeeds.

## Linked Issue

- #137
