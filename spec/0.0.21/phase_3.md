# v0.0.21 Phase 3 - CRD RBAC Sample Updates

## 작업 내용

- [x] Added CRD fields for `adminAuthProxyEnabled`, `adminAuthProxyImage`, `adminAuthProxyCookieSecretRef`, and `adminRouteHost`.
- [x] Added Role permissions for Jobs, Routes, and ServiceAccounts.
- [x] Updated the sample `OpsMateConfig`.

## 검증

- `kubectl kustomize config/crd`
- `go test ./...`
- `go build -o .cache\manager.exe ./cmd/manager`

## 남은 범위

- v0.0.22 should add install/runbook checks for cookie Secret creation and admin group binding.

## 연결 이슈

- #102
