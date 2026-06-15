# v0.0.32 Planner

## 1. Goal

- Upgrade the local CRC OLM install from the existing v0.0.29 CSV path to the published v0.0.31 catalog entry.
- Use the v0.0.31 manager logger to identify the first reconcile blocker.
- Prepare a v0.0.32 manager, bundle, and catalog fix for the discovered namespace cache/RBAC mismatch.

## 2. Architecture Overview

- CatalogSource stays in `openshift-marketplace`.
- Local OLM install stays isolated in `cywell-opsmate-olm`.
- The existing direct-bootstrap namespace remains untouched.
- v0.0.31 is the observed upgrade target because it includes manager zap logger initialization and updated bundle/catalog references.
- v0.0.32 configures manager cache scope from `WATCH_NAMESPACE` or `POD_NAMESPACE`, matching OwnNamespace/SingleNamespace OLM RBAC.

## 3. Technical Stack

- OpenShift CLI.
- OLM CatalogSource, Subscription, InstallPlan, CSV.
- `deploy/olm/local-crc` kustomize overlay.
- controller-runtime manager logs.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Refresh CatalogSource to v0.0.31 | done | CatalogSource READY evidence |
| Phase 2 | Approve and verify OLM upgrade | done | CSV v0.0.31 evidence |
| Phase 3 | Diagnose reconcile blocker and prepare v0.0.32 fix | done | namespace cache fix |
| Phase 4 | Document next handoff | done | v0.0.33 scope |

## 5. Migration Or Operation Strategy

- Do not delete local OLM resources unless upgrade recovery requires it and the destructive action is explicitly planned.
- Prefer in-place OLM upgrade from the existing Subscription.
- If the manager logs a namespace/RBAC mismatch, fix manager cache scoping before broadening RBAC.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| CatalogSource | `openshift-marketplace/cywell-opsmate-catalog` |
| Namespace | `cywell-opsmate-olm` |
| Observed upgrade target | `cywell-opsmate.v0.0.31` |
| Prepared upgrade target | `cywell-opsmate.v0.0.32` |
| Smoke CR | `cywell-opsmate-olm/cyops` |

## 7. Security Considerations

- Local CRC evidence is not production evidence.
- Do not commit kubeconfig, tokens, pull secrets, or unsanitized screenshots.
- The local PostgreSQL password Secret is a CRC-only smoke-test credential.

## 8. Completion Criteria

- [x] CatalogSource refreshed to the v0.0.31 catalog image and reported READY.
- [x] OLM upgrade reached CSV `cywell-opsmate.v0.0.31` `Succeeded` after a local recovery image patch.
- [x] Local `OpsMateConfig` reconcile result was captured with initialized logs.
- [x] v0.0.32 namespace-scoped cache fix and image/catalog references are prepared.
- [x] Follow-up work is documented and linked to GitHub Issues.
