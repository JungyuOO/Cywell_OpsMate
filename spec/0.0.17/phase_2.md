# Phase 2 - Migration Job Boundary

## Scope

- [x] Add a pgvector migration Job manifest builder.
- [x] Require `pgVectorMigrationApproved=true`.
- [x] Require DSN Secret reference and embedding dimensions.
- [x] Keep the Job out of reconciliation.
- [x] Link GitHub Issue #81.

## Work Completed

- Added `postgres.PGVectorMigrationJob`.
- Added `postgres.PGVectorMigrationJobName`.
- Added tests for approval gating and Secret-backed env.

## Verification

- `go test ./internal/controller/postgres`

## Remaining Scope

- Add the actual migration command image/entrypoint and admin trigger in a later version.
