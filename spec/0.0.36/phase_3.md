# v0.0.36 Phase 3 - Packaging Validation Cleanup

## Work

- [x] Added appserver routes for `/plugin-manifest.json`, `/console-plugin/plugin-manifest.json`, and `/console-plugin/plugin-entry.js`.
- [x] Added appserver tests for the new manifest and entry endpoints.
- [x] Updated manager, appserver, bundle, catalog, and CatalogSource version references to v0.0.36.
- [x] Preserved historical catalog relatedImages for older bundles.
- [x] Ran Go tests.
- [x] Attempted OpenShift/OLM manifest validation.
- [x] Recorded image-build validation blockers.

## Verification

- `go test ./...` passed.
- `rg -n "v0\.0\.36|v0\.0\.34|v0\.0\.27|0\.0\.36" .github deploy config internal\controller internal\appserver -S` confirmed current v0.0.36 references and historical appserver v0.0.27 references only in older catalog entries.
- `kubectl kustomize config/default` and `kubectl kustomize deploy/olm/local-crc` failed in this Codex sandbox with `evalsymlink failure ... Access is denied`.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` could not complete after CRC API became unreachable at `127.0.0.1:6443`.
- Docker-based local image build validation was not run because Docker Desktop API was unavailable at `npipe:////./pipe/dockerDesktopLinuxEngine`.

## Remaining Scope

- Re-run kustomize, OLM dry-run, local image builds, and CRC upgrade after Docker/CRC are available.
- The manifest/entry endpoints provide a loadable surface for the next smoke, but a full OpenShift dynamic plugin frontend bundle is still needed before CYOps can own a rich launcher/chat drawer experience.

## Linked Issue

- #177
