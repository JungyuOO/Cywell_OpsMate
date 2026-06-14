# Phase 2 - pgvector Live RAG Smoke

## Scope

- [x] Add a live smoke test gated by `CYOPS_PGVECTOR_TEST_DSN`.
- [x] Apply base migrations and controlled pgvector migration.
- [x] Seed ready document/chunk/vector rows.
- [x] Call `/api/chat` with pgvector retrieval mode.
- [x] Verify the first citation is the pgvector-ranked chunk.
- [x] Verify retrieval metrics recorded pgvector activity.
- [x] Link GitHub Issue #66.

## Work Completed

- Added `TestPGVectorLiveRAGSmoke`.
- The test verifies migration, SQL ranking, chat citation output, and metrics snapshot in one flow.

## Verification

- `CYOPS_PGVECTOR_TEST_DSN=postgres://cyops:cyops@localhost:55433/cyops?sslmode=disable go test ./internal/appserver -run TestPGVectorLiveRAGSmoke -count=1 -v`

## Remaining Scope

- Run the same test against an OpenShift-hosted pgvector database in v0.0.15.
