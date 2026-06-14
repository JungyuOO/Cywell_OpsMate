# Phase 3 - pgvector Rollout Boundary

## Scope

- [x] Keep bytea fallback as the dev/test path.
- [x] Keep `CYOPS_RETRIEVAL_MODE=pgvector` as the pgvector activation switch.
- [x] Document that OpenShift live validation requires a pgvector-enabled database image and controlled migration application.
- [x] Link GitHub Issue #62.

## Work Completed

- No schema change was added in this version.
- v0.0.13 focused on observability required before pgvector live rollout.
- Live validation remains blocked on a pgvector-enabled OpenShift database target.

## Verification

- Existing pgvector SQL builder and retrieval mode tests remain in the v0.0.12 path.
- v0.0.13 verification focuses on the observer-to-metrics path.

## Remaining Scope

- Apply `PGVectorEmbeddingMigrationSQL(dimensions)` against a pgvector-enabled database.
- Run `/api/chat` RAG smoke with `CYOPS_PGVECTOR_REQUIRED=true` and `CYOPS_RETRIEVAL_MODE=pgvector`.
