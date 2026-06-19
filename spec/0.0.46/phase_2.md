# v0.0.46 Phase 2 - Packaging And CRC Smoke

## Tasks

- [x] Bump manager, appserver, bundle, and catalog references to v0.0.46.
- [x] Keep catalog graph as `v0.0.46 -> v0.0.45 -> v0.0.44`.
- [x] Publish v0.0.46 images.
- [x] Upgrade CRC and verify plugin/provider content through gateway.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- GitHub Actions publish workflows succeeded for manager, appserver, bundle, and catalog on main.
- CRC CSV `cywell-opsmate.v0.0.46` reached `Succeeded`.
- CRC deployments use `ghcr.io/jungyuoo/cywell-opsmate:v0.0.46` and `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.46`.
- `oc rollout status deploy/console -n openshift-console` completed successfully.
- Gateway pod smoke confirmed `/plugin-entry.js` serves v0.0.46 and includes `X-CSRFToken`.

## Remaining Scope

- Authenticated browser chat confirmation remains user-observed after hard refresh.
