# v0.0.41 Phase 2 - Packaging And CRC Smoke

## 작업 항목

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.41.
- [x] Keep catalog graph as `v0.0.41 -> v0.0.40 -> v0.0.39`.
- [x] Publish v0.0.41 images.
- [x] Upgrade CRC and verify launcher entry bundle content.

## 검증

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- GitHub Actions publish workflows succeeded for manager, appserver, bundle, and catalog on main.
- CRC CSV `cywell-opsmate.v0.0.41` reached `Succeeded`.
- CRC deployments use `ghcr.io/jungyuoo/cywell-opsmate:v0.0.41` and `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.41`.
- Port-forward smoke for `/plugin-entry.js` and `/plugin-manifest.json` confirmed v0.0.41, callback registration, `CYOps`, `right: "22px"`, `bottom: "22px"`, and `2147483647`.

## 남은 범위

- Authenticated browser visual confirmation remains user-observed because Codex does not have the current OpenShift Console browser session credentials.
