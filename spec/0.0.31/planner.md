# v0.0.31 Planner

## 1. Goal

- Validate the CYOps OLM path locally on CRC before any real-server deployment.
- Avoid conflicts with the existing development bootstrap namespace.
- Capture the exact commands and evidence for CatalogSource, Subscription, InstallPlan, CSV, and first `OpsMateConfig` reconcile checks.
- Prepare a v0.0.31 manager, bundle, and catalog revision because local OLM smoke exposed that the existing manager image can crash without useful controller-runtime logs.

## 2. Architecture Overview

- CatalogSource remains cluster-level in `openshift-marketplace`.
- Local OLM Subscription uses isolated namespace `cywell-opsmate-olm`.
- Existing direct-bootstrap namespace `cywell-opsmate-system` is left untouched.
- InstallPlan approval remains manual.
- The local overlay includes a sanitized `OpsMateConfig` and PostgreSQL password Secret for CRC-only reconcile testing.
- v0.0.31 catalog upgrades from `cywell-opsmate.v0.0.29` to `cywell-opsmate.v0.0.31`.

## 3. Technical Stack

- OpenShift CLI.
- OLM CatalogSource, OperatorGroup, Subscription, InstallPlan, CSV.
- Kustomize overlay for local CRC install.
- Go manager logging via controller-runtime zap logger.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Isolated local OLM manifests | done | `deploy/olm/local-crc` |
| Phase 2 | CRC catalog verification | done | CatalogSource READY evidence |
| Phase 3 | CRC Subscription, InstallPlan, CSV smoke | done | CSV `Succeeded` evidence |
| Phase 4 | OpsMateConfig smoke finding and v0.0.31 upgrade handoff | done | logging fix, version bump, v0.0.32 handoff |

## 5. Migration Or Operation Strategy

- Do not delete development bootstrap resources during this version.
- Install local OLM resources into `cywell-opsmate-olm`.
- Approve InstallPlan manually only after the generated plan is reviewed.
- Treat the v0.0.29 CSV success as OLM-path evidence and the `OpsMateConfig` manager crash as the next upgrade validation target.
- Publish v0.0.31 manager, bundle, and catalog images before retrying local reconcile on CRC.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| CatalogSource | `openshift-marketplace/cywell-opsmate-catalog` |
| Local OLM namespace | `cywell-opsmate-olm` |
| Subscription | `cywell-opsmate` |
| InstallPlan approval | Manual |
| Local smoke CR | `cywell-opsmate-olm/cyops` |

## 7. Security Considerations

- Local CRC evidence is not real-server evidence.
- Do not commit kubeconfig, tokens, pull secrets, or unsanitized screenshots.
- Existing direct-bootstrap resources are preserved unless cleanup is explicitly planned.
- The committed local PostgreSQL Secret uses a CRC-only test password and must not be reused as a customer credential.

## 8. Completion Criteria

- [x] Local CRC OLM manifests exist.
- [x] Manifests validate with client dry-run.
- [x] CatalogSource reached READY on CRC.
- [x] Subscription generated an InstallPlan and CSV reached `Succeeded`.
- [x] Local `OpsMateConfig` smoke exposed the manager logging gap.
- [x] v0.0.31 manager, bundle, catalog references are prepared for publish.
- [x] Next upgrade and Web Console evidence handoff is documented.
