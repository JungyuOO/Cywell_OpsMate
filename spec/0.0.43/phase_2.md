# v0.0.43 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.43.
- [x] Keep catalog graph as `v0.0.43 -> v0.0.42 -> v0.0.41`.
- [x] Publish v0.0.43 images.
- [x] Upgrade CRC and verify plugin entry content through gateway.

## Verification

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- GitHub Actions publish workflows succeeded for manager, appserver, bundle, and catalog on main.
- CRC CSV `cywell-opsmate.v0.0.43` reached `Succeeded`.
- CRC deployments use `ghcr.io/jungyuoo/cywell-opsmate:v0.0.43` and `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.43`.
- `oc rollout status deploy/console -n openshift-console` completed successfully.
- Gateway port-forward smoke confirmed v0.0.43, `cyops-console@0.0.43`, `cyopsLauncherFlag`, `console.flag`, `right: "22px"`, and `2147483647`.

## Remaining Scope

- Authenticated browser visual confirmation remains user-observed because Codex does not have the current OpenShift Console browser session credentials.
