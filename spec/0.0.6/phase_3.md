# Phase 3 - storage/provider wiring refinements

## 작업 내용

- [x] appserver dependency wiring 방식을 문서화한다.
- [x] storage adapter와 provider client를 config에서 주입할 경계를 정리한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] `go test ./...`

## 남은 범위

- [ ] 실제 Lightspeed REST API client는 후속 버전에서 완성한다.
