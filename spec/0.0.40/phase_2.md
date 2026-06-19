# v0.0.40 Phase 2 - Packaging And CRC Smoke

## 작업 항목

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.40.
- [x] Keep catalog graph as `v0.0.40 -> v0.0.39 -> v0.0.38`.
- [ ] Publish v0.0.40 images.
- [ ] Upgrade CRC and verify launcher entry bundle content.

## 검증

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.

## 남은 범위

- GitHub Issue and PR handoff remain for Phase 3.
