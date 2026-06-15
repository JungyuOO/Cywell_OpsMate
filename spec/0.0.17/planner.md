# v0.0.17 Planner

## 1. Goal

- Add an explicit re-embedding batch workflow for ready documents.
- Add a pgvector migration Job manifest boundary that refuses to build without admin approval.
- Keep migration Job creation out of reconciliation.
- Document the next runtime status and console diagnostics work.

## 2. Architecture Overview

1. `EmbeddingService.ReembedReadyDocuments` selects ready documents and re-runs the existing embedding path.
2. Failed re-embedding uses the existing `FailEmbedding` path so `embedding_status=failed` and `last_error` are observable.
3. `postgres.PGVectorMigrationJob` builds a Secret-backed migration Job template only when `pgVectorMigrationApproved=true`.
4. Reconciliation still does not create the migration Job automatically.

## 3. Tech Stack

| Area | Choice | Reason |
| --- | --- | --- |
| Re-embedding | appserver service method | reuses existing embedding provider and repository behavior |
| Document selection | PostgreSQL query | avoids loading non-ready/deleted documents |
| Migration boundary | Kubernetes Job manifest builder | prepares OpenShift execution without enabling automatic migration |
| Approval gate | `pgVectorMigrationApproved` | keeps admin approval explicit |

## 4. Implementation Steps

| Phase | Scope | Status | Notes |
| --- | --- | --- | --- |
| Phase 1 | re-embedding batch workflow | done | ready document selection and failure status |
| Phase 2 | migration Job boundary | done | approval, DSN Secret, dimensions required |
| Phase 3 | verification | done | unit and Postgres integration tests |
| Phase 4 | v0.0.18 handoff | done | runtime status and console diagnostics next |

## Linked Issues

- Phase 1: #80
- Phase 2: #81
- Phase 3: #82
- Phase 4: #83

## 5. Migration Or Operation Strategy

- Migration Job builder does not run unless explicitly called by future admin automation.
- The Job template uses `CYOPS_POSTGRES_DSN` from a SecretKeyRef and never embeds the DSN value.
- Re-embedding is callable after migration to regenerate ready documents.
- Failed re-embedding leaves document state visible through `embedding_status` and `last_error`.

## 6. Message / Communication / Data Protocol

| Path | Payload | Status |
| --- | --- | --- |
| Repository -> Reembedding | ready documents with embedding status pending/ready/failed | done |
| Reembedding -> Repository | `BeginEmbedding`, `CompleteEmbedding`, `FailEmbedding` | existing |
| CR -> Job builder | approval, DSN SecretRef, dimensions | done |
| Job env | `CYOPS_POSTGRES_DSN`, `CYOPS_EMBEDDING_DIMENSIONS` | done |

## 7. Security Considerations

- Migration Job template references DSN through SecretKeyRef only.
- Tests verify the DSN value is not placed directly in the Job env.
- Re-embedding errors are stored as operational `last_error`; provider implementations must continue avoiding token/DSN leakage.
- Reconciliation still avoids silent destructive migration.

## 8. Completion Criteria

- [x] Re-embedding can be triggered explicitly from appserver code.
- [x] Failed re-embedding sets `embedding_status=failed` and `last_error`.
- [x] Migration Job builder requires approval before returning a manifest.
- [x] Migration Job template uses SecretKeyRef for DSN.
- [x] v0.0.18 handoff captures runtime status and console diagnostics.
