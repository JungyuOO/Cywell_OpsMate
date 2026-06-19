# v0.0.42 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.42.
- [x] Keep catalog graph as `v0.0.42 -> v0.0.41 -> v0.0.40`.
- [ ] Publish v0.0.42 images.
- [ ] Upgrade CRC and verify gateway-backed plugin endpoint content.

## Verification

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.

## Remaining Scope

- GitHub Issue and PR handoff remain for Phase 3.
