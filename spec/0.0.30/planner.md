# v0.0.30 Planner

## 1. Goal

- Add a catalog image path for local/CRC OLM testing before any real-server deployment.
- Add a CatalogSource manifest that points at the CYOps catalog image.
- Fix cluster-scoped ConsolePlugin RBAC for both direct bootstrap and OLM installs.

## 2. Architecture Overview

- Bundle image `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.29` contains the CSV and CRD.
- Catalog image `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.30` exposes the bundle through OLM file-based catalog metadata.
- `CatalogSource` lives in `openshift-marketplace`.
- `OperatorGroup` and `Subscription` live in `cywell-opsmate-system`.
- ConsolePlugin permissions are cluster-scoped through ClusterRole/clusterPermissions.

## 3. Technical Stack

- OLM file-based catalog.
- `opm` catalog image runtime.
- GitHub Actions GHCR publishing.
- OpenShift client dry-run validation.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Catalog image metadata | done | `deploy/olm/catalog/catalog.yaml` |
| Phase 2 | Catalog image publish workflow | done | catalog Containerfile + workflow |
| Phase 3 | Local OLM install manifests | done | CatalogSource + install docs |
| Phase 4 | RBAC correction and local handoff | done | ClusterRole/CSV clusterPermissions |

## 5. Migration Or Operation Strategy

- Use direct `config/default` apply only for development bootstrap.
- Use OLM CatalogSource, OperatorGroup, and Subscription for local/CRC installation tests.
- Only after CRC OLM install is validated should the same catalog path be pointed at a real server.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| Catalog image | `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.30` |
| CatalogSource namespace | `openshift-marketplace` |
| Subscription namespace | `cywell-opsmate-system` |
| Subscription approval | Manual InstallPlan approval |
| ConsolePlugin RBAC | ClusterRole/clusterPermissions |

## 7. Security Considerations

- Catalog and bundle images contain manifests only, not tokens.
- Manual InstallPlan approval is retained while the install path is still being tested.
- CRC/local testing must not be confused with real customer-server validation.

## 8. Completion Criteria

- [x] Catalog metadata exists.
- [x] Catalog image workflow exists.
- [x] CatalogSource manifest exists.
- [x] ConsolePlugin RBAC is cluster-scoped.
- [x] Local validations pass.
