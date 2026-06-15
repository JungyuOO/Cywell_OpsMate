# Phase 1 - Database Spec And pgvector Image

## Scope

- [x] Add database image configuration to `OpsMateConfig`.
- [x] Allow an explicit custom image.
- [x] Default managed PostgreSQL to a pgvector-enabled image when `requirePGVector` is true.
- [x] Link GitHub Issue #70.

## Work Completed

- Added `DatabaseSpec.Image`.
- Added `postgres.imageFor`.
- Added controller tests for explicit image and pgvector-required default image.

## Verification

- `go test ./internal/controller/postgres`

## Remaining Scope

- Production image pinning policy should be finalized in the deployment runbook.
