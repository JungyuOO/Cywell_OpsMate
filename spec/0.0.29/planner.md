# v0.0.29 Planner

## 1. Goal

- Start replacing direct `oc apply` installation with an OLM-native Operator distribution path.
- Add a CYOps Operator bundle with CSV, CRD, metadata, and bundle image publishing.
- Provide install manifests that model the customer-facing OperatorHub flow.

## 2. Architecture Overview

- Bundle image: `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.29`.
- CSV installs the manager Deployment using `ghcr.io/jungyuoo/cywell-opsmate:v0.0.28`.
- CSV declares the `OpsMateConfig` owned CRD and required namespace-scoped permissions.
- Customers should install through CatalogSource + OperatorGroup + Subscription rather than raw kustomize apply.

## 3. Technical Stack

- OLM ClusterServiceVersion.
- OLM registry bundle metadata.
- GitHub Actions GHCR bundle image publish.
- OpenShift OLM install manifests.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | OLM bundle manifests | done | CSV + CRD + metadata |
| Phase 2 | Bundle image publishing | done | bundle Containerfile + workflow |
| Phase 3 | OLM install manifests | done | namespace, OperatorGroup, Subscription |
| Phase 4 | Catalog handoff | done | v0.0.30 handoff |

## 5. Migration Or Operation Strategy

- Keep `config/default` for development bootstrap and controller smoke tests.
- Use OLM bundle/catalog for realistic OpenShift Operator delivery.
- Use manual InstallPlan approval until the install path is tested end-to-end.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| Bundle image | `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.29` |
| Package | `cywell-opsmate` |
| Channel | `alpha` |
| Subscription source | `cywell-opsmate-catalog` |
| Subscription namespace | `openshift-marketplace` |

## 7. Security Considerations

- Bundle metadata contains image references, not secrets.
- InstallPlan approval is manual for the initial OLM path.
- Evidence captured from Web Console or network tools must be sanitized.

## 8. Completion Criteria

- [x] OLM bundle manifests exist.
- [x] Bundle image workflow exists.
- [x] OLM install manifests document OperatorGroup/Subscription usage.
- [x] Local validation passes.
