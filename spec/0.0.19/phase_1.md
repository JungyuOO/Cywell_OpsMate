# Phase 1 - Admin Endpoint Authorization

## Scope

- [x] Add `CYOPS_ADMIN_TOKEN`.
- [x] Inject token through appserver Deployment SecretKeyRef.
- [x] Require `X-CYOps-Admin-Token` for `/api/ops/reembed`.
- [x] Link GitHub Issue #90.

## Work Completed

- Added `AdminToken` config and server option.
- Added `ConsoleSpec.AdminTokenSecretRef` and `AdminTokenSecretKey`.
- Added tests for missing admin token.

## Verification

- `go test ./internal/appserver ./internal/controller/appserver ./api/v1alpha1`
- `CYOPS_POSTGRES_TEST_DSN=postgres://cyops:cyops@localhost:55432/cyops?sslmode=disable go test ./internal/appserver -count=1`

## Remaining Scope

- Replace token header with OpenShift OAuth/RBAC integration before production exposure.
