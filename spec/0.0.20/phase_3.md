# v0.0.20 Phase 3 - CRD Sample and Verification

## 작업 내용

- [x] Added `spec.console.adminUsers`.
- [x] Added `spec.console.adminGroups`.
- [x] Passed admin users/groups into the appserver Deployment.
- [x] Updated CRD schema and sample manifest.

## 검증

- `go test ./internal/appserver ./internal/controller ./internal/controller/appserver ./internal/controller/postgres ./api/v1alpha1`
- `kubectl kustomize config/crd`
- `go test ./...`
- `go build -o .cache\manager.exe ./cmd/manager`
- `go build -o .cache\cyops-pgvector-migrate.exe ./cmd/cyops-pgvector-migrate`

## 남은 범위

- Console admin/debug panel placement is deferred until the OpenShift route/oauth-proxy surface exists.

## 연결 이슈

- #97
