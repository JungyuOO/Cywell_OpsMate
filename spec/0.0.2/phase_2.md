# Phase 2 - PostgreSQL desired resource helper

## 작업 내용

- [x] `internal/controller/postgres`에 Deployment/Service builder를 추가한다.
- [x] PostgreSQL image, DB 이름, 포트, 기본 설정을 명시한다.
- [x] PostgreSQL Secret/PVC 자동 생성은 구현하지 않고 후속 범위로 기록한다.
- [x] desired object shape를 단위 테스트로 고정한다.

## 검증

- [x] `go test ./internal/controller/postgres`
- [x] `go test ./...`

## 남은 범위

- [ ] password Secret 생성/참조는 후속 버전에서 구현한다.
- [ ] PVC, backup, HA 구성은 후속 버전에서 구현한다.
