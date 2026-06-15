# Phase 2 - Appserver DSN Secret Wiring

## Scope

- [x] Add database DSN Secret fields to `OpsMateConfig`.
- [x] Inject `CYOPS_POSTGRES_DSN` into appserver from SecretKeyRef.
- [x] Default DSN key to `dsn` when omitted.
- [x] Link GitHub Issue #71.

## Work Completed

- Added `DatabaseSpec.DSNSecretRef`.
- Added `DatabaseSpec.DSNSecretKey`.
- Updated appserver Deployment env generation and tests.

## Verification

- `go test ./internal/controller/appserver`

## Remaining Scope

- The operator still expects the Secret to exist; automatic Secret creation for managed Postgres is a later decision.
