# v0.0.32 Phase 3 - Reconcile Blocker And Namespace Cache Fix

## Work

- [x] Collected initialized v0.0.31 manager logs.
- [x] Identified that manager cache was listing `OpsMateConfig` at cluster scope.
- [x] Kept namespace-scoped OLM RBAC and fixed the manager to watch only `WATCH_NAMESPACE` or `POD_NAMESPACE`.
- [x] Added `POD_NAMESPACE` Downward API env to direct and OLM manager deployments.
- [x] Bumped manager, bundle, and catalog references to v0.0.32.

## Verification

- `oc logs deployment/cywell-opsmate-controller-manager -n cywell-opsmate-olm --tail=120`
- `$env:GOCACHE = (Join-Path (Get-Location) '.cache\go-build'); go test ./...`
- `kubectl kustomize config/default`
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests`
- `oc apply --dry-run=client -k deploy\olm\local-crc`
- `docker build -f deploy\containerfiles\manager.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate:v0.0.32 .`
- `docker build -f deploy\containerfiles\bundle.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.32 .`
- `docker build -f deploy\containerfiles\catalog.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.32 .`
- `docker run --rm -v ${PWD}\deploy\olm\catalog:/configs:ro quay.io/operator-framework/opm:latest validate /configs`

Observed blocker:

- `failed to list *v1alpha1.OpsMateConfig: opsmateconfigs.opsmate.cywell.io is forbidden ... cannot list resource "opsmateconfigs" ... at the cluster scope`

## Remaining Scope

- Publish v0.0.32 images, refresh CRC CatalogSource, approve the v0.0.32 upgrade, and retry reconcile.

## Linked Issue

- #156
