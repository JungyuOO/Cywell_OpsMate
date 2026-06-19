# v0.0.53 Phase 3 - Packaging And CRC Smoke

## Tasks

- [ ] Create and link GitHub issue.
- [ ] Publish v0.0.53 manager, appserver, bundle, and catalog images.
- [ ] Upgrade CRC to `cywell-opsmate.v0.0.53`.
- [ ] Confirm path-based chat smoke returns a Lightspeed response.
- [ ] Confirm plugin entry no longer includes drawer document upload UI.

## Verification

- Created #216 for v0.0.53 center document workspace and path chat fallback.
- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.

## Remaining Scope

- Complete after v0.0.53 is installed on CRC.
