# Phase 3 - Storage 및 Lightspeed provider skeleton

## 작업 내용

- [x] document storage adapter interface를 추가한다.
- [x] local/PVC storage skeleton을 추가한다.
- [x] Lightspeed provider interface와 mocked implementation을 분리한다.
- [x] customer document context 최소 전달 정책을 코드 주석/테스트로 고정한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] `go test ./...`

## 남은 범위

- [ ] 실제 Lightspeed REST API 호출은 v0.0.6 이후로 이관한다.
