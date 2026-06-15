# v0.0.20 Phase 1 - OpenShift Admin Identity Authorization

## 작업 내용

- [x] Added `AdminAuthConfig` to appserver options.
- [x] Preserved `X-CYOps-Admin-Token` authorization.
- [x] Added `X-Forwarded-User` allowlist authorization.
- [x] Added CSV `X-Forwarded-Groups` allowlist authorization.
- [x] Added env config for `CYOPS_ADMIN_USERS` and `CYOPS_ADMIN_GROUPS`.

## 검증

- `go test ./internal/appserver`

## 남은 범위

- v0.0.21 should add the OpenShift route/oauth-proxy deployment shape that makes forwarded headers trustworthy in production.

## 연결 이슈

- #95
