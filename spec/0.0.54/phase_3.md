# v0.0.54 Phase 3 - Packaging And CRC Smoke

## Tasks

- [x] Create and link GitHub issue.
- [x] Publish v0.0.54 manager, appserver, bundle, and catalog images.
- [x] Upgrade CRC to `cywell-opsmate.v0.0.54`.
- [x] Confirm plugin assets report `0.0.54`.
- [x] Confirm the Documents page is served as a page-style surface.

## Verification

- Created #218 for v0.0.54 OpenShift Console-style Documents page alignment.
- PR #219 merged and closed #218.
- GitHub Actions for PR #219 completed successfully: `operator-catalog`, `operator-bundle`, `manager-image`, and `appserver-image`.
- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- CRC CatalogSource applied with `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.54`.
- CRC InstallPlan `install-l4fr5` was approved.
- CRC CSV check passed: `cywell-opsmate.v0.0.54` is `Succeeded` and replaces `cywell-opsmate.v0.0.53`.
- CRC deployments are ready: `cyops-appserver`, `cyops-gateway`, `cyops-postgres`, and `cywell-opsmate-controller-manager` are all `1/1`.
- CRC plugin smoke passed:
  - `/plugin-manifest.json` reports version `0.0.54`.
  - `/console-plugin/documents` contains `cyops-page-header`, `cyops-toolbar`, `cyops-table`, and `자료 관리`.
  - `/console-plugin/plugin-entry.js` reports version `0.0.54` and navigates the injected `자료` action to `/console-plugin/documents`.

## Remaining Scope

- Browser visual verification remains manual because the CRC console certificate warning blocks automated in-app browser navigation without bypassing the warning.
