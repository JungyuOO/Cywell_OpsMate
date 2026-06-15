# v0.0.27 Planner

## 1. Goal

- Publish the appserver image for OpenShift live testing.
- Pin the generated appserver and migration workloads to the published image tag.
- Capture initial live-cluster readiness evidence without changing the Red Hat Lightspeed viewer.

## 2. Architecture Overview

- `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.27` is the appserver image used by the Operator-generated Deployment.
- The pgvector migration Job uses the same image tag and runs `cyops-pgvector-migrate`.
- The Web Console path remains same-origin through the ConsolePlugin backend service.

## 3. Technical Stack

- Docker local build and GitHub Actions GHCR publish.
- Go controller desired object builders.
- OpenShift CLI read/apply verification.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Version-pinned appserver image references | done | controller image defaults |
| Phase 2 | GHCR image publish workflow | done | `.github/workflows/appserver-image.yml` |
| Phase 3 | OpenShift readiness probe | done | live CLI evidence |
| Phase 4 | Next live UI evidence handoff | done | v0.0.28 handoff |

## 5. Migration Or Operation Strategy

- Build the appserver image locally before merge, then let GitHub Actions publish it with repository-scoped `packages: write`.
- Keep the migration command in the appserver image so migration Jobs do not need a separate image family.
- Keep `latest` out of generated appserver workloads to make test-cluster evidence reproducible.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| Appserver Deployment image | `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.27` |
| Migration Job image | `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.27` |
| Migration command | `cyops-pgvector-migrate` |
| ConsolePlugin backend path | same-origin Web Console backend service |

## 7. Security Considerations

- Local GHCR login uses tokens only as process input and must not print tokens.
- The publish workflow uses the GitHub Actions `GITHUB_TOKEN` with `packages: write`; local personal tokens are not required to carry package scope.
- Live evidence must be sanitized before committing.
- The Red Hat OpenShift Lightspeed viewer remains untouched.

## 8. Completion Criteria

- [x] Controller appserver and migration image defaults are pinned to `v0.0.27`.
- [x] The appserver image is built locally and a GHCR publish workflow is added.
- [x] OpenShift CLI access/readiness is verified.
- [x] Remaining screenshot/network evidence is handed off.
