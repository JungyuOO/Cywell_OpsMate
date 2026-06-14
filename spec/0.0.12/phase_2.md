# Phase 2 - pgvector SQL ranking path

## 작업 항목

- [x] pgvector top-k SQL builder를 추가한다.
- [x] repository query method boundary를 추가한다.
- [x] pgvector mode가 SQL ranking path를 선택하도록 한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- pgvector-enabled live DB 실행 검증은 후속 버전에서 진행한다.

## 작업 내용

- `rankedChunksPGVectorSQL`과 `ListRankedChunksPGVector`를 추가했다.
- `PostgresRetriever`가 `pgvector` mode에서 SQL ranking path를 선택한다.
