# v0.0.48 Phase 3 - Packaging And Handoff

## Tasks

- [x] Create or link GitHub Issues.
- [x] Commit with Lore trailers.
- [x] Push branch and open PR.
- [x] Merge PR and verify linked Issues close.
- [x] Publish and install v0.0.48 on CRC.

## Verification

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- Created #206 for the v0.0.48 Lightspeed boundary fix.
- Commit `07a457a` was pushed on `feature/v0.0.48`.
- PR #207 was merged into `main`; issue #206 is closed.
- GitHub Actions completed successfully for `manager-image`, `appserver-image`, `operator-bundle`, and `operator-catalog`.
- CRC OLM upgraded to `cywell-opsmate.v0.0.48` with phase `Succeeded`.
- CRC deployments use `ghcr.io/jungyuoo/cywell-opsmate:v0.0.48` and `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.48`.
- Appserver env smoke confirmed no `CYOPS_LIGHTSPEED_MODEL`, `CYOPS_LIGHTSPEED_PROVIDER`, `LIGHTSPEED_DEFAULT_MODEL`, or `LIGHTSPEED_DEFAULT_PROVIDER`.
- Gateway smoke confirmed `plugin-entry.js` version `0.0.48`.

## Remaining Scope

- Apply the Lightspeed `OLSConfig` example on a real cluster with Lightspeed Operator installed.
