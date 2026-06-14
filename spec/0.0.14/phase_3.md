# Phase 3 - Local pgvector Verification

## Scope

- [x] Start a local `pgvector/pgvector:pg16` database.
- [x] Execute the env-gated live smoke.
- [x] Confirm existing BYTEA integration tests still pass against the non-pgvector test database.
- [x] Link GitHub Issue #67.

## Work Completed

- Local pgvector smoke passed against port 55433.
- Existing `CYOPS_POSTGRES_TEST_DSN` integration suite passed against the regular PostgreSQL test container.
- Temporary pgvector test container was stopped and removed.

## Verification

- `CYOPS_POSTGRES_TEST_DSN=postgres://cyops:cyops@localhost:55432/cyops?sslmode=disable go test ./internal/appserver -count=1`
- `CYOPS_PGVECTOR_TEST_DSN=postgres://cyops:cyops@localhost:55433/cyops?sslmode=disable go test ./internal/appserver -run TestPGVectorLiveRAGSmoke -count=1 -v`

## Remaining Scope

- OpenShift live validation still needs cluster-hosted pgvector database credentials.
