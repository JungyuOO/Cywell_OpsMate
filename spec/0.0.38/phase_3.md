# v0.0.38 Phase 3 - Packaging And CRC Smoke

## 작업 항목

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.38.
- [ ] Publish v0.0.38 images through GitHub Actions.
- [ ] Upgrade CRC through OLM to `cywell-opsmate.v0.0.38`.
- [ ] Verify the manifest, entry script, and launcher path from CRC.

## 검증

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- Catalog graph now has `v0.0.38 -> v0.0.37 -> v0.0.34`.

## 남은 범위

- Issue closure, PR merge, and final handoff remain for Phase 4.
