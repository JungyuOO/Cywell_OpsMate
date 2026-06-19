# v0.0.49 Phase 3 - Packaging And CRC Smoke

## Tasks

- [x] Create or link GitHub Issues.
- [ ] Commit with Lore trailers.
- [ ] Push branch and open PR.
- [ ] Merge PR and verify linked Issues close.
- [ ] Publish and install v0.0.49 on CRC.
- [ ] Patch CYOps `OpsMateConfig` to Lightspeed `/v1/query`.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- Created #208 for the v0.0.49 Lightspeed OLS query wiring.

## Remaining Scope

- Complete packaging and CRC smoke.
