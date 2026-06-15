# Phase 2 - Runtime Status Fields And Conditions

## Scope

- [x] Add `status.pgVectorReady`.
- [x] Add aggregate `status.reembedding`.
- [x] Reflect runtime fields through conditions.
- [x] Degrade status when re-embedding has failures.
- [x] Link GitHub Issue #91.

## Work Completed

- Added `ReembeddingStatus`.
- Updated `PGVectorReady` condition to use runtime evidence.
- Added `ReembeddingReady` condition.
- Added controller tests.

## Verification

- `go test ./internal/controller ./internal/controller/appserver ./api/v1alpha1`

## Remaining Scope

- Wire actual migration Job/re-embedding workflows to update these status fields.
