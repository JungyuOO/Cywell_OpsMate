# v0.0.52 Phase 3 - Packaging And CRC Smoke

## Tasks

- [ ] Create and link GitHub issue.
- [ ] Publish v0.0.52 manager, appserver, bundle, and catalog images.
- [ ] Upgrade CRC to `cywell-opsmate.v0.0.52`.
- [ ] Confirm backend chat smoke still returns Lightspeed response.
- [ ] Confirm browser chat no longer shows `/api/chat returned 403`.

## Verification

- Created #214 for v0.0.52 Console chat input and document navigation.
- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.

## Remaining Scope

- Complete after v0.0.52 is installed on CRC.
