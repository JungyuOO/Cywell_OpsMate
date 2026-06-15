# Phase 2 - Operator Readiness Conditions

## Scope

- [x] Add `PostgresDSNConfigured`.
- [x] Add `PGVectorRequired`.
- [x] Add `PGVectorMigrationApproved`.
- [x] Add `PGVectorReady`.
- [x] Add `RetrievalModeReady`.
- [x] Link GitHub Issue #76.

## Work Completed

- Added status condition builders to the reconciler.
- Added tests for default config, valid pgvector config, and invalid pgvector retrieval mode.
- `PGVectorReady` is `Unknown` when runtime appserver/live smoke validation is pending.

## Verification

- `go test ./internal/controller`

## Remaining Scope

- Add real runtime status updates after migration job and live DB checks are implemented.
