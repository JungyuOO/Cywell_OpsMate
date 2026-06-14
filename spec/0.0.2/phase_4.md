# Phase 4 - 검증 및 v0.0.3 이관 범위 정리

## 작업 내용

- [ ] v0.0.2 phase 완료 상태를 `planner.md`와 phase 문서에 반영한다.
- [ ] Go 검증과 manifest 검증 결과를 기록한다.
- [ ] v0.0.3으로 넘길 ConsolePlugin, backend API, bundle/catalog 범위를 정리한다.

## 검증

- [ ] `go fmt ./...`
- [ ] `go test ./...`
- [ ] `go build -o .cache/manager.exe ./cmd/manager`
- [ ] `kubectl kustomize config/default`

## 남은 범위

- [ ] OpenShift ConsolePlugin reconcile
- [ ] Lightspeed backend API 구현
- [ ] appserver TLS/service-ca 연결
- [ ] PostgreSQL Secret/PVC/HA 운영화
- [ ] Operator bundle/catalog 정리
