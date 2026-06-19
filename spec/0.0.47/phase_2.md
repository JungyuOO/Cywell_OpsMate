# v0.0.47 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.47.
- [x] Keep catalog graph as `v0.0.47 -> v0.0.46 -> v0.0.45`.
- [ ] Publish v0.0.47 images.
- [ ] Upgrade CRC and verify plugin/provider content through gateway.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.

## Remaining Scope

- Merge release PR and wait for image workflows.
- Upgrade CRC after catalog publication.
