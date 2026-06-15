# v0.0.23 Phase 3 - Fallback Route Opt-in Cleanup

## 작업 내용

- [x] Changed sample `adminAuthProxyEnabled` to `false`.
- [x] Removed default sample `adminRouteHost`.
- [x] Updated OpenShift smoke script to test fallback Route only with `-EnableFallbackRoute`.

## 검증

- PowerShell parser validation for `deploy/scripts/openshift-v022-smoke.ps1`.
- `kubectl kustomize config/crd`.

## 남은 범위

- Live fallback Route smoke remains available for direct admin access scenarios.

## 연결 이슈

- #112
