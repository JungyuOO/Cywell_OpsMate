# v0.0.18 Planner

## 1. Goal

- Add the real `cyops-pgvector-migrate` command entrypoint used by the migration Job boundary.
- Add an admin-triggered appserver re-embedding endpoint with batch limit support.
- Keep responses and logs free of DSNs, tokens, prompts, filenames, and document text.
- Document the next runtime status and console diagnostics work.

## 2. Architecture Overview

1. `cmd/cyops-pgvector-migrate` reads `CYOPS_POSTGRES_DSN` and `CYOPS_EMBEDDING_DIMENSIONS`.
2. The command calls `ApplyPGVectorEmbeddingMigration`.
3. `POST /api/ops/reembed` invokes `EmbeddingService.ReembedReadyDocuments` for Postgres-backed appservers.
4. The endpoint returns only processed and failed counts.

## 3. Tech Stack

| Area | Choice | Reason |
| --- | --- | --- |
| Migration entrypoint | Go CLI command | matches existing Go runtime and Job builder |
| Re-embedding trigger | appserver admin API | immediate admin workflow without adding another controller surface |
| Batch control | request `limit` | supports incremental re-embedding |
| Security | Secret env + count-only response | avoids exposing customer content or credentials |

## 4. Implementation Steps

| Phase | Scope | Status | Notes |
| --- | --- | --- | --- |
| Phase 1 | migration CLI entrypoint | done | `cyops-pgvector-migrate` |
| Phase 2 | admin re-embedding API | done | `POST /api/ops/reembed` |
| Phase 3 | verification | done | unit, integration, full test, builds |
| Phase 4 | v0.0.19 handoff | done | runtime status and console diagnostics next |

## Linked Issues

- Phase 1: #85
- Phase 2: #86
- Phase 3: #87
- Phase 4: #88

## 5. Migration Or Operation Strategy

- Migration execution remains admin-triggered through the migration Job boundary.
- The command fails fast when DSN or dimensions are missing.
- Re-embedding runs in batches and reports aggregate counts.
- OpenShift RBAC/auth policy for `/api/ops/reembed` must be finalized before production exposure.

## 6. Message / Communication / Data Protocol

| Path | Payload | Status |
| --- | --- | --- |
| Job env -> CLI | `CYOPS_POSTGRES_DSN`, `CYOPS_EMBEDDING_DIMENSIONS` | done |
| CLI -> PostgreSQL | controlled pgvector migration | done |
| Admin -> appserver | `POST /api/ops/reembed` with optional `limit` | done |
| appserver -> Admin | `processed`, `failed` counts | done |

## 7. Security Considerations

- CLI validation errors do not print DSN values.
- Endpoint response does not include document IDs, filenames, chunk text, prompts, or provider details.
- Endpoint currently relies on the appserver/OpenShift route auth boundary; production RBAC remains follow-up.
- Provider errors may be stored in `last_error`, so provider adapters must continue avoiding token leakage.

## 8. Completion Criteria

- [x] Migration Job has a real executable entrypoint.
- [x] CLI config is env-driven and tested.
- [x] Re-embedding can be triggered by an admin API.
- [x] Re-embedding endpoint returns aggregate counts only.
- [x] Full tests and both binaries build.
