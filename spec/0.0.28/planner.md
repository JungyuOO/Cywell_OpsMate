# v0.0.28 Planner

## 1. Goal

- Add a publishable Operator manager image.
- Pin OpenShift manager deployment to `ghcr.io/jungyuoo/cywell-opsmate:v0.0.28`.
- Apply the merged manifests to OpenShift after the manager image is published.

## 2. Architecture Overview

- The manager image runs `/manager`.
- GitHub Actions publishes `ghcr.io/jungyuoo/cywell-opsmate:v0.0.28` and `latest`.
- `config/default` applies namespace, CRD, RBAC, and manager Deployment.
- The manager then reconciles `OpsMateConfig` resources into appserver/PostgreSQL/ConsolePlugin resources.

## 3. Technical Stack

- Go controller manager.
- Docker multi-stage manager image build.
- GitHub Actions GHCR publishing with `packages: write`.
- OpenShift CLI apply/readiness checks.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Manager image packaging | done | `deploy/containerfiles/manager.Containerfile` |
| Phase 2 | Manager image publish workflow | done | `.github/workflows/manager-image.yml` |
| Phase 3 | OpenShift manifest image pin | done | `config/manager/manager.yaml` |
| Phase 4 | Apply handoff | done | `v0.0.29_handoff.md` |

## 5. Migration Or Operation Strategy

- Merge the manager image workflow first.
- Wait for `manager-image` to publish `v0.0.28`.
- Apply `kubectl kustomize config/default | oc apply -f -`.
- Create or update an `OpsMateConfig` only after the manager pod is available.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| Manager Deployment image | `ghcr.io/jungyuoo/cywell-opsmate:v0.0.28` |
| Manager entrypoint | `/manager` |
| Apply command | `kubectl kustomize config/default | oc apply -f -` |
| Sample CR command | `oc apply -f config/samples/opsmate_v1alpha1_opsmateconfig.yaml` |

## 7. Security Considerations

- GHCR publishing uses GitHub Actions `GITHUB_TOKEN` with `packages: write`.
- The manager image runs as non-root in the Deployment.
- Do not commit kubeconfig, OpenShift tokens, GHCR tokens, screenshots with sensitive cluster identifiers, or raw network headers.

## 8. Completion Criteria

- [x] Manager image Containerfile exists.
- [x] Manager image workflow exists.
- [x] Manager Deployment points at the versioned image.
- [x] Tests and local image build pass.
