# v0.0.50 Phase 2 - Packaging

## Tasks

- [x] Bump manager, appserver, bundle, catalog, and plugin versions to v0.0.50.
- [x] Keep catalog graph as `v0.0.50 -> v0.0.49 -> v0.0.48`.
- [ ] Publish and install v0.0.50 on CRC.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.

## Remaining Scope

- Run tests, merge PR, wait for images, and approve CRC InstallPlan.
