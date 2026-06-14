# Phase 3 - reconciler create/update 및 status

## 작업 내용

- [x] `OpsMateConfigReconciler`가 CR을 조회한다.
- [x] appserver와 PostgreSQL desired object에 owner reference를 설정한다.
- [x] Deployment/Service를 create/update한다.
- [x] `status.conditions`와 `overallStatus`를 최소 갱신한다.
- [x] RBAC에 apps/core/status 권한을 추가한다.

## 검증

- [x] controller fake client 테스트
- [x] `go test ./...`
- [x] `go build -o .cache/manager.exe ./cmd/manager`
- [x] `kubectl kustomize config/default`

## 남은 범위

- [ ] ConsolePlugin reconcile은 v0.0.3으로 이관한다.
- [ ] Lightspeed backend API 구현은 v0.0.3 이후로 이관한다.
