# v0.0.42 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.42.
- [x] Keep catalog graph as `v0.0.42 -> v0.0.41 -> v0.0.40`.
- [x] Publish v0.0.42 images.
- [x] Upgrade CRC and verify gateway-backed plugin endpoint content.

## Verification

- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- GitHub Actions publish workflows succeeded for manager, appserver, bundle, and catalog on main.
- CRC CSV `cywell-opsmate.v0.0.42` reached `Succeeded`.
- CRC deployments include `cyops-gateway` running `nginxinc/nginx-unprivileged:1.27-alpine`.
- ConsolePlugin `cyops-console` backend points to Service `cyops-gateway` on port `8443`.
- Gateway port-forward smoke confirmed `/plugin-entry.js` and `/plugin-manifest.json` include v0.0.42, callback registration, `CYOps`, `right: "22px"`, `bottom: "22px"`, and `2147483647`.
- `oc rollout status deploy/console -n openshift-console` completed successfully.

## Remaining Scope

- Authenticated browser visual confirmation remains user-observed because Codex does not have the current OpenShift Console browser session credentials.
