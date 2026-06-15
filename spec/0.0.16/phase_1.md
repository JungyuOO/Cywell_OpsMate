# Phase 1 - Migration Approval Field

## Scope

- [x] Add explicit admin approval field.
- [x] Keep default value false.
- [x] Avoid automatic migration execution in reconciliation.
- [x] Link GitHub Issue #75.

## Work Completed

- Added `DatabaseSpec.PGVectorMigrationApproved`.
- Added `PGVectorMigrationApproved` status condition.

## Verification

- `go test ./internal/controller`

## Remaining Scope

- Future migration automation must check this field before creating a job.
