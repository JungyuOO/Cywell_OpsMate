# Phase 2 - Admin Re-embedding API

## Scope

- [x] Add `POST /api/ops/reembed`.
- [x] Accept optional batch `limit`.
- [x] Require a Postgres-backed document repository.
- [x] Return processed and failed counts only.
- [x] Link GitHub Issue #86.

## Work Completed

- Added `ReembeddingAPIRequest`.
- Added `ReembeddingAPIResponse`.
- Wired appserver embedder into server options.
- Added endpoint tests and Postgres integration coverage.

## Verification

- `go test ./internal/appserver`
- `CYOPS_POSTGRES_TEST_DSN=postgres://cyops:cyops@localhost:55432/cyops?sslmode=disable go test ./internal/appserver -count=1`

## Remaining Scope

- Add production RBAC/auth policy for admin-only access.
