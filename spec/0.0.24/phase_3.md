# v0.0.24 Phase 3 - Verification and No-extra-OAuth Checks

## 작업 내용

- [x] Added tests for the served diagnostics HTML.
- [x] Added tests that diagnostics JS calls `/api/ops/diagnostics` and schema.
- [x] Added tests that the console diagnostics JS does not handle OAuth route redirects.

## 검증

- `go test ./internal/appserver ./internal/controller/console`
- `go test ./...`
- `kubectl kustomize config/crd`

## 남은 범위

- Browser-level verification in OpenShift Web Console remains a live-cluster follow-up.

## 연결 이슈

- #117
