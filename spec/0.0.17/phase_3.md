# Phase 3 - Verification

## Scope

- [x] Run appserver tests.
- [x] Run controller postgres tests.
- [x] Run Postgres integration tests with Docker.
- [x] Run full repository tests and build.
- [x] Link GitHub Issue #82.

## Work Completed

- Re-embedding failure path was verified against PostgreSQL.
- Migration Job boundary tests verify approval and SecretKeyRef behavior.

## Verification

- `go test ./internal/appserver ./internal/controller/postgres`
- `CYOPS_POSTGRES_TEST_DSN=postgres://cyops:cyops@localhost:55432/cyops?sslmode=disable go test ./internal/appserver -count=1`
- `go test ./...`
- `go build -o .cache\manager.exe ./cmd/manager`

## Remaining Scope

- OpenShift Job execution is not enabled yet.
