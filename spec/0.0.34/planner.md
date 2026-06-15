# v0.0.34 Planner

## 1. Goal

- Apply the published v0.0.33 catalog to CRC.
- Upgrade the local OLM install to `cywell-opsmate.v0.0.33`.
- Verify the ConsolePlugin update/resourceVersion fix under live reconcile.
- Capture the next blocker after ConsolePlugin reconciliation is clean and prepare the v0.0.34 fix.

## 2. Architecture Overview

- CatalogSource remains in `openshift-marketplace`.
- Subscription remains in `cywell-opsmate-olm`.
- v0.0.33 manager should reconcile cluster-scoped ConsolePlugin resources without dropping server metadata.
- v0.0.34 keeps Service annotations during reconcile so OpenShift serving cert Secrets are requested, and mounts writable PostgreSQL data/run directories for restricted SCC.
- Existing direct-bootstrap namespace remains untouched.

## 3. Technical Stack

- OpenShift CLI.
- OLM CatalogSource, InstallPlan, CSV, Deployment.
- `deploy/olm/local-crc` overlay.
- controller-runtime logs.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Refresh CatalogSource to v0.0.33 | done | CatalogSource evidence |
| Phase 2 | Approve and verify v0.0.33 upgrade | done | CSV/deployment evidence |
| Phase 3 | Retry `OpsMateConfig` reconcile and fix remaining pod blockers | done | resource/log evidence, v0.0.34 fix |
| Phase 4 | Handoff | done | v0.0.35 scope |

## 5. Migration Or Operation Strategy

- Prefer in-place OLM upgrade from the current local Subscription.
- Do not delete local OLM resources unless explicitly required by a documented recovery step.
- If v0.0.32 crash loops block OLM progress, use a local recovery image patch and document it as CRC-only.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| CatalogSource | `openshift-marketplace/cywell-opsmate-catalog` |
| Namespace | `cywell-opsmate-olm` |
| Upgrade target | `cywell-opsmate.v0.0.33` |
| Prepared fix target | `cywell-opsmate.v0.0.34` |
| Smoke CR | `cywell-opsmate-olm/cyops` |

## 7. Security Considerations

- CRC evidence is local validation only.
- Do not commit kubeconfig, tokens, pull secrets, or unsanitized screenshots.
- CRC-only Secret values must not be reused as customer credentials.

## 8. Completion Criteria

- [x] CatalogSource refreshed to v0.0.33 and reported READY.
- [x] CSV `cywell-opsmate.v0.0.33` reached `Succeeded`.
- [x] Manager ran with v0.0.33 image.
- [x] ConsolePlugin reconcile no longer logged `resourceVersion` update errors.
- [x] Remaining appserver/PostgreSQL blockers were captured with evidence.
- [x] v0.0.34 Service annotation and PostgreSQL writable directory fix is prepared.
