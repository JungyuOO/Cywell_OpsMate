# v0.0.33 Planner

## 1. Goal

- Apply the published v0.0.32 catalog to CRC.
- Upgrade the local OLM install to `cywell-opsmate.v0.0.32`.
- Verify the namespace-scoped manager cache removes the cluster-scope `OpsMateConfig` RBAC error.
- Retry the local `OpsMateConfig/cyops` reconcile, capture the next blocker, and prepare the v0.0.33 fix.

## 2. Architecture Overview

- CatalogSource remains in `openshift-marketplace`.
- Subscription remains in `cywell-opsmate-olm`.
- v0.0.32 manager should watch only the pod namespace through `POD_NAMESPACE`.
- v0.0.33 preserves existing metadata/resourceVersion when reconciling unstructured resources such as ConsolePlugin.
- Existing direct-bootstrap namespace remains untouched.

## 3. Technical Stack

- OpenShift CLI.
- OLM CatalogSource, InstallPlan, CSV, Deployment.
- `deploy/olm/local-crc` overlay.
- controller-runtime logs.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Refresh CatalogSource to v0.0.32 | done | CatalogSource evidence |
| Phase 2 | Approve and verify v0.0.32 upgrade | done | CSV/deployment evidence |
| Phase 3 | Retry `OpsMateConfig` reconcile and fix ConsolePlugin update | done | resource/log evidence, v0.0.33 fix |
| Phase 4 | Handoff | done | v0.0.34 scope |

## 5. Migration Or Operation Strategy

- Prefer in-place OLM upgrade from the current local Subscription.
- Do not delete local OLM resources unless explicitly required by a documented recovery step.
- If new blockers appear, capture initialized logs before code changes.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| CatalogSource | `openshift-marketplace/cywell-opsmate-catalog` |
| Namespace | `cywell-opsmate-olm` |
| Upgrade target | `cywell-opsmate.v0.0.32` |
| Prepared fix target | `cywell-opsmate.v0.0.33` |
| Smoke CR | `cywell-opsmate-olm/cyops` |

## 7. Security Considerations

- CRC evidence is local validation only.
- Do not commit kubeconfig, tokens, pull secrets, or unsanitized screenshots.
- CRC-only Secret values must not be reused as customer credentials.

## 8. Completion Criteria

- [x] CatalogSource refreshed to v0.0.32 and reported READY.
- [x] CSV `cywell-opsmate.v0.0.32` reached `Succeeded`.
- [x] Manager ran with v0.0.32 image and no cluster-scope `OpsMateConfig` RBAC errors.
- [x] `OpsMateConfig/cyops` reconcile result was captured.
- [x] ConsolePlugin update/resourceVersion fix is prepared as v0.0.33.
- [x] GitHub Issues and phase docs are linked.
