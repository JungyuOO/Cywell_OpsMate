# Phase 3 - ranked retrieval and citations

## 작업 항목

- [x] stored embeddings를 이용해 retrieval score를 계산한다.
- [x] top-k ranking을 provider context에 반영한다.
- [x] citations에 rank/score를 포함한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] Docker PostgreSQL integration test

## 남은 범위

- pgvector SQL similarity ranking은 후속 migration 뒤 구현한다.

## 작업 내용

- `PostgresRetriever`가 query embedding과 stored BYTEA embeddings를 cosine score로 rank한다.
- citations에 `rank`와 `score`를 추가했다.
