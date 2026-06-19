# v0.0.51 Phase 2 - Packaging

## Tasks

- [x] Bump manager, appserver, bundle, catalog, and plugin versions to v0.0.51.
- [x] Keep catalog graph as `v0.0.51 -> v0.0.50 -> v0.0.49`.
- [x] Publish and install v0.0.51 on CRC.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- PR #213 merged to `main`.
- GitHub Actions succeeded for `appserver-image`, `manager-image`, `operator-bundle`, and `operator-catalog`.
- CRC CSV `cywell-opsmate.v0.0.51` reached `Succeeded`.

## Remaining Scope

- None for packaging.
