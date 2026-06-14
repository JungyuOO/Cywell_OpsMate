# Phase 4 - v0.0.7 ingestion handoff

## 작업 내용

- [x] v0.0.6 완료 상태를 문서에 반영한다.
- [x] v0.0.7 ingestion/embedding 구현 범위를 정리한다.

## 검증

- [x] `go fmt ./...`
- [x] `go test ./...`
- [x] `go build -o .cache/manager.exe ./cmd/manager`

## 남은 범위

- [ ] parser/chunker implementation
- [ ] embedding worker
- [ ] pgvector migration update
- [ ] frontend API integration
