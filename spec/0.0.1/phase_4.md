# Phase 4 - 검증 및 v0.0.2 이관 범위 정리

## 작업 내용

- [ ] v0.0.1 phase 완료 상태를 `planner.md`와 phase 문서에 반영한다.
- [ ] Go 검증과 manifest 검증 결과를 기록한다.
- [ ] v0.0.2로 넘길 실제 reconcile 구현 범위를 정리한다.

## 검증

- [ ] `go fmt ./...`
- [ ] `go test ./...`
- [ ] `go build -o .cache/manager.exe ./cmd/manager`
- [ ] manifest 검증 명령

## 남은 범위

- [ ] appserver Deployment/Service reconcile
- [ ] PostgreSQL 리소스 reconcile
- [ ] OpenShift ConsolePlugin reconcile
- [ ] Lightspeed API client와 backend API 구현
