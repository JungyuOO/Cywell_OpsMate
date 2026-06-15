# v0.0.19 Planner

## 1. Goal

- Add admin authorization for `/api/ops/reembed`.
- Add runtime status fields for pgvector readiness and re-embedding progress.
- Wire admin token through OpenShift Secret-backed appserver env.
- Document OpenShift migration Job execution validation as the next required live step.

## 2. Architecture Overview

1. `OpsMateConfig.spec.console.adminTokenSecretRef/key` provides `CYOPS_ADMIN_TOKEN`.
2. `POST /api/ops/reembed` requires `X-CYOps-Admin-Token`.
3. `OpsMateConfig.status.pgVectorReady` can mark runtime pgvector validation success.
4. `OpsMateConfig.status.reembedding` carries aggregate re-embedding state without document details.
5. Conditions reflect runtime fields through `PGVectorReady` and `ReembeddingReady`.

## 3. Tech Stack

| Area | Choice | Reason |
| --- | --- | --- |
| Admin auth | Secret-backed bearer-style header | simple boundary until OpenShift OAuth/RBAC integration |
| Status | CR status fields + conditions | visible in OpenShift and controller tests |
| Re-embedding response | aggregate counts | avoids customer metadata leakage |
| Live validation | deferred OpenShift Job runbook | requires target cluster execution |

## 4. Implementation Steps

| Phase | Scope | Status | Notes |
| --- | --- | --- | --- |
| Phase 1 | admin endpoint authorization | done | token env + header check |
| Phase 2 | runtime status fields and conditions | done | pgvector and re-embedding |
| Phase 3 | CRD/sample verification | done | schema and env SecretRef |
| Phase 4 | v0.0.20 handoff | done | OpenShift Job validation next |

## Linked Issues

- Phase 1: #90
- Phase 2: #91
- Phase 3: #92
- Phase 4: #93

## 5. Migration Or Operation Strategy

- Re-embedding admin endpoint is disabled unless `CYOPS_ADMIN_TOKEN` is configured.
- The admin token is injected from a Secret through `OpsMateConfig.spec.console`.
- Runtime status fields are aggregate-only and must not contain DSNs, tokens, prompts, filenames, or document text.
- Migration Job execution still needs target OpenShift validation.

## 6. Message / Communication / Data Protocol

| Path | Payload | Status |
| --- | --- | --- |
| Console Secret -> appserver | `CYOPS_ADMIN_TOKEN` | done |
| Admin -> appserver | `X-CYOps-Admin-Token` | done |
| Runtime -> CR status | `pgVectorReady` | done |
| Runtime -> CR status | `reembedding.running/processed/failed/lastError` | done |
| CR status -> conditions | `PGVectorReady`, `ReembeddingReady` | done |

## 7. Security Considerations

- Missing or mismatched admin token returns 403.
- Re-embedding API returns only processed and failed counts.
- CR status does not include secrets or document content.
- This is an interim admin auth boundary; OpenShift OAuth/RBAC integration remains a follow-up.

## 8. Completion Criteria

- [x] `/api/ops/reembed` requires admin authorization.
- [x] Appserver can receive admin token from Secret-backed env.
- [x] Runtime pgvector readiness can be reflected in status.
- [x] Re-embedding progress can be reflected in status and conditions.
- [x] CRD/sample/tests/builds are updated.
