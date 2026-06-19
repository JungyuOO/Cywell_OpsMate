# v0.0.53 Phase 3 - Packaging And CRC Smoke

## Tasks

- [x] Create and link GitHub issue.
- [x] Publish v0.0.53 manager, appserver, bundle, and catalog images.
- [x] Upgrade CRC to `cywell-opsmate.v0.0.53`.
- [x] Confirm path-based chat smoke returns a Lightspeed response.
- [x] Confirm plugin entry no longer includes drawer document upload UI.

## Verification

- Created #216 for v0.0.53 center document workspace and path chat fallback.
- PR #217 merged and closed #216.
- GitHub Actions for PR #217 completed successfully: `operator-catalog`, `operator-bundle`, `manager-image`, and `appserver-image`.
- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- CRC CSV check passed: `cywell-opsmate.v0.0.53` is `Succeeded`.
- CRC deployments are ready: `cyops-appserver`, `cyops-gateway`, `cyops-postgres`, and `cywell-opsmate-controller-manager` are all `1/1`.
- CRC backend smoke passed through the service path endpoint:
  - `GET /api/chat/message/Say%20hello%20from%20CYOps`
  - response was HTTP 200 with provider `lightspeed`.
- Plugin asset smoke confirmed `plugin-manifest.json` version `0.0.53` includes the Documents navigation entry and `plugin-entry.js` no longer renders the chat drawer document upload/list panel.

## Remaining Scope

- Browser visual verification remains manual because the CRC console certificate warning blocks automated in-app browser navigation without bypassing the warning.
- Future versions should replace the injected DOM workspace with official OpenShift Console page extension APIs when the plugin page surface is expanded.
