# v0.0.54 Phase 3 - Packaging And CRC Smoke

## Tasks

- [x] Create and link GitHub issue.
- [ ] Publish v0.0.54 manager, appserver, bundle, and catalog images.
- [ ] Upgrade CRC to `cywell-opsmate.v0.0.54`.
- [ ] Confirm plugin assets report `0.0.54`.
- [ ] Confirm the Documents page is served as a page-style surface.

## Verification

- Created #218 for v0.0.54 OpenShift Console-style Documents page alignment.
- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.

## Remaining Scope

- Complete after PR and CRC upgrade.
