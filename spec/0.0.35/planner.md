# v0.0.35 Planner

## 1. Goal

- Apply the published v0.0.34 catalog to CRC.
- Upgrade the local OLM install to `cywell-opsmate.v0.0.34`.
- Verify the v0.0.34 serving-cert Service annotation and PostgreSQL restricted-SCC fixes in a live CRC install.
- Capture any remaining readiness blocker with enough evidence for the next version.

## 2. Architecture Overview

- CatalogSource remains in `openshift-marketplace`.
- Subscription remains in `cywell-opsmate-olm`.
- OLM should perform an in-place upgrade from the current local CSV.
- The manager should reconcile `OpsMateConfig` into appserver/PostgreSQL resources that can become ready under CRC restricted SCC.

## 3. Technical Stack

- OpenShift CLI.
- CRC OpenShift cluster with OLM.
- OLM CatalogSource, InstallPlan, CSV, Deployment, Service, Secret, and Pod resources.
- `deploy/olm/install/catalogsource.yaml`.
- controller-runtime and workload pod logs.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Confirm v0.0.34 release artifacts | done | GHCR and workflow evidence |
| Phase 2 | Refresh CRC CatalogSource to v0.0.34 | done | CatalogSource and catalog pod evidence |
| Phase 3 | Approve and verify v0.0.34 OLM upgrade | done | InstallPlan/CSV/deployment evidence |
| Phase 4 | Verify appserver/PostgreSQL readiness | done | pod, secret, service, and log evidence |
| Phase 5 | Handoff | done | v0.0.36 scope |

## 5. Migration Or Operation Strategy

- Prefer in-place OLM upgrade from the current local Subscription.
- Do not delete local OLM resources unless explicitly required by a documented recovery step.
- Re-trigger the existing smoke `OpsMateConfig` reconcile with an annotation update after the upgrade.
- Treat CRC-only patches as diagnostic recovery steps and document them separately from product behavior.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| CatalogSource | `openshift-marketplace/cywell-opsmate-catalog` |
| Namespace | `cywell-opsmate-olm` |
| Upgrade target | `cywell-opsmate.v0.0.34` |
| Smoke CR | `cywell-opsmate-olm/cyops` |
| Appserver service | `cyops-appserver` |
| PostgreSQL deployment | `cyops-postgres` |

## 7. Security Considerations

- CRC evidence is local validation only.
- Do not commit kubeconfig, tokens, pull secrets, generated TLS keys, or unsanitized credentials.
- Do not broaden SCC or add privileged permissions unless a later version explicitly scopes that change.

## 8. Completion Criteria

- [x] v0.0.34 workflows and GHCR images are confirmed.
- [x] CatalogSource refreshed to v0.0.34 and reported READY.
- [x] CSV `cywell-opsmate.v0.0.34` reached `Succeeded`.
- [x] Manager ran with v0.0.34 image.
- [x] Appserver serving cert Secret was created or the missing Secret blocker was re-captured.
- [x] PostgreSQL pod became ready or the new restricted-SCC failure was captured.
- [x] Remaining scope is assigned to the next version.
