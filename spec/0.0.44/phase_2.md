# v0.0.44 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.44.
- [x] Keep catalog graph as `v0.0.44 -> v0.0.43 -> v0.0.42`.
- [x] Publish v0.0.44 images.
- [x] Upgrade CRC and verify plugin entry content through gateway.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- `opm` was not available in the local PATH, so catalog validation is limited to YAML review and OpenShift dry-run.
- GitHub Actions publish workflows succeeded for manager, appserver, bundle, and catalog on main.
- CRC CSV `cywell-opsmate.v0.0.44` reached `Succeeded`.
- CRC deployments use `ghcr.io/jungyuoo/cywell-opsmate:v0.0.44` and `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.44`.
- `oc rollout status deploy/console -n openshift-console` completed successfully.
- Gateway pod smoke confirmed `/plugin-entry.js` and `/plugin-manifest.json` serve v0.0.44, `cyops-console@0.0.44`, `data-cyops-plugin-entry`, and delayed mount retry content.

## Remaining Scope

- Authenticated browser visual confirmation remains user-observed because Codex does not have the current OpenShift Console browser session credentials.
