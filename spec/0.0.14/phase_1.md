# Phase 1 - Controlled pgvector Migration Helper

## Scope

- [x] Add a helper that validates dimensions before database mutation.
- [x] Check pgvector extension readiness before applying the generated migration SQL.
- [x] Preserve default BYTEA migration behavior.
- [x] Link GitHub Issue #65.

## Work Completed

- Added `ApplyPGVectorEmbeddingMigration`.
- Added a regression test that ensures invalid dimensions fail before DB execution.

## Verification

- `go test ./internal/appserver`

## Remaining Scope

- OpenShift operator wiring for invoking the controlled migration remains a follow-up.
