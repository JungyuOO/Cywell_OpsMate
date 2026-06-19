# v0.0.38 Phase 3 - Packaging And CRC Smoke

## 작업 항목

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.38.
- [x] Publish v0.0.38 images through GitHub Actions.
- [x] Upgrade CRC through OLM to `cywell-opsmate.v0.0.38`.
- [x] Verify the manifest, entry script, and launcher path from CRC.

## 검증

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- Catalog graph now has `v0.0.38 -> v0.0.37 -> v0.0.34`.
- GitHub Actions published v0.0.38 manager, appserver, bundle, and catalog images.
- `docker manifest inspect` succeeded for all four v0.0.38 GHCR tags.
- CRC `CatalogSource/cywell-opsmate-catalog` was applied with `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.38`.
- InstallPlan `install-57zc5` was approved and CSV `cywell-opsmate.v0.0.38` reached `Succeeded`.
- `cywell-opsmate-controller-manager` ran `ghcr.io/jungyuoo/cywell-opsmate:v0.0.38`.
- `cyops-appserver` rolled out to `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.38`.
- CRC port-forward smoke returned:
  - `manifestHasVersion: true`
  - `manifestHasCallback: true`
  - `manifestHasBaseURL: true`
  - `entryHasCallback: true`
  - `entryHasLauncher: true`
  - `entryHasChat: true`
- Direct Console proxy curl with bearer token still returned `401 Unauthorized`; authenticated browser-session smoke remains the right way to verify the Console proxy path.

## 남은 범위

- Authenticated OpenShift Web Console visual smoke can be repeated from a browser session that already trusts the CRC certificate and has a console cookie.
