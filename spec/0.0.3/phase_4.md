# Phase 4 - backend HTTP skeleton 및 검증

## 작업 내용

- [x] appserver backend HTTP skeleton과 health endpoint를 추가한다.
- [x] v0.0.3 phase 완료 상태를 `planner.md`와 phase 문서에 반영한다.
- [x] v0.0.4로 넘길 Lightspeed API client, frontend bundle, PostgreSQL 운영화 범위를 정리한다.

## 검증

- [x] `go fmt ./...`
- [x] `go test ./...`
- [x] `go build -o .cache/manager.exe ./cmd/manager`
- [x] `kubectl kustomize config/default`

## 남은 범위

- [ ] Lightspeed API client 구현
- [ ] ConsolePlugin frontend bundle
- [ ] PostgreSQL Secret/PVC/HA 운영화
- [ ] Operator bundle/catalog 정리
