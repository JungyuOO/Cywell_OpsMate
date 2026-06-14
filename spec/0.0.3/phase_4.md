# Phase 4 - backend HTTP skeleton 및 검증

## 작업 내용

- [ ] appserver backend HTTP skeleton과 health endpoint를 추가한다.
- [ ] v0.0.3 phase 완료 상태를 `planner.md`와 phase 문서에 반영한다.
- [ ] v0.0.4로 넘길 Lightspeed API client, frontend bundle, PostgreSQL 운영화 범위를 정리한다.

## 검증

- [ ] `go fmt ./...`
- [ ] `go test ./...`
- [ ] `go build -o .cache/manager.exe ./cmd/manager`
- [ ] `kubectl kustomize config/default`

## 남은 범위

- [ ] Lightspeed API client 구현
- [ ] ConsolePlugin frontend bundle
- [ ] PostgreSQL Secret/PVC/HA 운영화
- [ ] Operator bundle/catalog 정리
