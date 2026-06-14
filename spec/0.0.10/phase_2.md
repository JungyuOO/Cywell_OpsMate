# Phase 2 - external embedding provider

## 작업 항목

- [x] embedding endpoint/model/dimensions config를 추가한다.
- [x] HTTP embedding provider를 구현한다.
- [x] provider request/response contract를 테스트한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- provider auth Secret wiring은 Operator config 확장 시점에 보강한다.

## 작업 내용

- `HTTPEmbeddingProvider`를 추가했다.
- `CYOPS_EMBEDDING_ENDPOINT`, `CYOPS_EMBEDDING_MODEL`, `CYOPS_EMBEDDING_DIMENSIONS` config를 추가했다.
