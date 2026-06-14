# Phase 2 - ConsolePlugin reconcile 및 RBAC

## 작업 내용

- [x] `OpsMateConfigReconciler`가 ConsolePlugin desired object를 create/update한다.
- [x] `console.enabled=false`일 때 ConsolePlugin 생성을 건너뛴다.
- [x] RBAC에 `console.openshift.io` ConsolePlugin 권한을 추가한다.
- [x] fake client 테스트로 ConsolePlugin reconcile을 검증한다.

## 검증

- [x] `go test ./internal/controller`
- [x] `go test ./...`
- [x] `kubectl kustomize config/default`

## 남은 범위

- [ ] ConsolePlugin frontend bundle과 navigation extension은 후속 버전으로 이관한다.
