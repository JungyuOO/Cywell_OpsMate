# v0.0.47 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.47.
- [x] Keep catalog graph as `v0.0.47 -> v0.0.46 -> v0.0.45`.
- [x] Publish v0.0.47 images.
- [x] Upgrade CRC and verify plugin/provider content through gateway.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- GitHub Actions completed successfully for `manager-image`, `appserver-image`, `operator-bundle`, and `operator-catalog`.
- CRC OLM upgraded to `cywell-opsmate.v0.0.47` with phase `Succeeded`.
- CRC deployments use `ghcr.io/jungyuoo/cywell-opsmate:v0.0.47` and `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.47`.
- `oc rollout status deploy/console -n openshift-console --timeout=180s` passed.
- Gateway smoke confirmed `plugin-entry.js` version `0.0.47` and non-GET headers `X-CSRFToken`, `X-CSRF-Token`, and `X-Requested-With`.

## Remaining Scope

- Browser chat smoke remains to be repeated in the authenticated console tab.
