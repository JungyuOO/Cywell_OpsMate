# Phase 1 - appserver desired resource helper

## 작업 내용

- [x] `internal/controller/appserver`에 Deployment/Service builder를 추가한다.
- [x] Lightspeed API endpoint와 Secret 참조를 appserver 환경변수로 전달한다.
- [x] appserver image 기본값을 한 곳에서 관리한다.
- [x] desired object shape를 단위 테스트로 고정한다.

## 검증

- [x] `go test ./internal/controller/appserver`
- [x] `go test ./...`

## 남은 범위

- [ ] 실제 backend HTTP API 구현은 후속 버전으로 이관한다.
- [ ] TLS/service-ca 연결은 ConsolePlugin phase 또는 후속 버전으로 이관한다.
