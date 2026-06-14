# Phase 4 - 검증 및 v0.0.2 이관 범위 정리

## 작업 내용

- [x] v0.0.1 phase 완료 상태를 `planner.md`와 phase 문서에 반영한다.
- [x] Go 검증과 manifest 검증 결과를 기록한다.
- [x] v0.0.2로 넘길 실제 reconcile 구현 범위를 정리한다.

## 검증

- [x] `go fmt ./...`
- [x] `go test ./...`
- [x] `go build -o .cache/manager.exe ./cmd/manager`
- [x] `kubectl kustomize config/default`
- [x] `kubectl kustomize config/samples`

## 남은 범위

- [ ] appserver Deployment/Service reconcile
- [ ] PostgreSQL 리소스 reconcile
- [ ] OpenShift ConsolePlugin reconcile
- [ ] Lightspeed API client와 backend API 구현
- [ ] Operator bundle/catalog 정리
