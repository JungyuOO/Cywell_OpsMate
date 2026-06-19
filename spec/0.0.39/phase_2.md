# v0.0.39 Phase 2 - Packaging And CRC Smoke

## 작업 항목

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.39.
- [x] Keep the catalog graph as `v0.0.39 -> v0.0.38 -> v0.0.37`.
- [x] Publish v0.0.39 images.
- [x] Upgrade CRC and verify plugin endpoint shape.

## 검증

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- GitHub Actions published v0.0.39 manager, appserver, bundle, and catalog images.
- `docker manifest inspect` succeeded for all four v0.0.39 GHCR tags.
- CRC `CatalogSource/cywell-opsmate-catalog` was applied with v0.0.39.
- InstallPlan `install-d4m2l` was approved and CSV `cywell-opsmate.v0.0.39` reached `Succeeded`.
- `cywell-opsmate-controller-manager` ran `ghcr.io/jungyuoo/cywell-opsmate:v0.0.39`.
- `cyops-appserver` ran `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.39`.
- CRC port-forward endpoint smoke returned:
  - `version: 0.0.39`
  - `hasTopLevelDisplayName: false`
  - `hasTopLevelDescription: false`
  - `displayName: CYOps`
  - `registrationMethod: callback`
  - `loadScript: plugin-entry.js`
  - `entryHasCallback: true`
  - `entryHasLauncher: true`

## 남은 범위

- Authenticated Web Console browser reload is still required to visually confirm the launcher in the real console session.
