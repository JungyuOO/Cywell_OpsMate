# v0.0.44 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.44.
- [x] Keep catalog graph as `v0.0.44 -> v0.0.43 -> v0.0.42`.
- [ ] Publish v0.0.44 images.
- [ ] Upgrade CRC and verify plugin entry content through gateway.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- `opm` was not available in the local PATH, so catalog validation is limited to YAML review and OpenShift dry-run.

## Remaining Scope

- Merge release PR and wait for image workflows.
- Upgrade CRC after catalog publication.
