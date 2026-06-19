# v0.0.52 Phase 3 - Packaging And CRC Smoke

## Tasks

- [x] Create and link GitHub issue.
- [x] Publish v0.0.52 manager, appserver, bundle, and catalog images.
- [x] Upgrade CRC to `cywell-opsmate.v0.0.52`.
- [x] Confirm backend chat smoke still returns Lightspeed response.
- [ ] Confirm browser chat no longer shows `/api/chat returned 403`.

## Verification

- Created #214 for v0.0.52 Console chat input and document navigation.
- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- PR #215 merged and closed #214.
- GitHub Actions succeeded for `appserver-image`, `manager-image`, `operator-bundle`, and `operator-catalog`.
- CRC CSV `cywell-opsmate.v0.0.52` is `Succeeded`.
- CRC appserver Deployment image is `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.52`.
- CRC manager Deployment image is `ghcr.io/jungyuoo/cywell-opsmate:v0.0.52`.
- CRC `POST /api/chat` smoke from `deploy/cyops-gateway` returned HTTP 200 with provider `lightspeed`.
- CRC `GET /api/chat?message=...&provider=lightspeed&rag=true` smoke from `deploy/cyops-gateway` returned HTTP 200 with provider `lightspeed`.
- CRC `/console-plugin/documents` returned HTTP 200.
- CRC `plugin-manifest.json` reports version `0.0.52` and includes `cyops-documents`.
- CRC `plugin-entry.js` reports version `0.0.52` and includes `requestChat`, Enter handling, and `cyops-console-nav-documents`.
- Automated browser smoke was not completed because the browser automation session stopped at the CRC self-signed certificate warning.

## Remaining Scope

- Manually hard-refresh the OpenShift Console tab and confirm the visual chat no longer shows `/api/chat returned 403`.
