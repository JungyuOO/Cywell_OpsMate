# Phase 1 - retrieval mode config

## 작업 항목

- [x] `CYOPS_RETRIEVAL_MODE` config를 추가한다.
- [x] appserver config가 retriever mode를 주입한다.
- [x] 기본값은 `bytea` fallback으로 유지한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- pgvector SQL ranking은 Phase 2에서 구현했다.

## 작업 내용

- appserver config에 `CYOPS_RETRIEVAL_MODE`와 `CYOPS_RETRIEVAL_SLOW_THRESHOLD_MS`를 추가했다.
- `OpsMateConfig.spec.embedding`에도 retrieval mode와 slow threshold boundary를 추가했다.
