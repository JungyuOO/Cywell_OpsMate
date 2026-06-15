# Phase 1 - Migration CLI Entrypoint

## Scope

- [x] Add `cmd/cyops-pgvector-migrate`.
- [x] Load DSN and dimensions from env.
- [x] Call `ApplyPGVectorEmbeddingMigration`.
- [x] Avoid printing DSN values in validation errors.
- [x] Link GitHub Issue #85.

## Work Completed

- Added `PGVectorMigrationCommandConfig`.
- Added `PGVectorMigrationCommandConfigFromEnv`.
- Added `RunPGVectorMigrationCommand`.
- Added CLI `main`.

## Verification

- `go test ./internal/appserver ./internal/controller/postgres ./cmd/cyops-pgvector-migrate`
- `go build -o .cache\cyops-pgvector-migrate.exe ./cmd/cyops-pgvector-migrate`

## Remaining Scope

- OpenShift Job execution against a real cluster remains follow-up.
