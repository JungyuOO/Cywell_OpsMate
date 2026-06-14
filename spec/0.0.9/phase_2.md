# Phase 2 - embedding persistence

## 작업 항목

- [x] chunk embeddings를 교체 저장하는 repository method를 추가한다.
- [x] embedding status를 processing/ready/failed로 갱신한다.
- [x] Docker PostgreSQL integration test를 추가한다.

## 검증

- [x] Docker PostgreSQL integration test

## 남은 범위

- pgvector `VECTOR(n)` migration은 v0.0.10 이후로 넘긴다.

## 작업 내용

- `BeginEmbedding`, `CompleteEmbedding`, `FailEmbedding`, `ReplaceEmbeddings`, `ListEmbeddingsContext`를 추가했다.
- `EmbeddingService`가 chunks를 조회하고 embedding rows를 저장하도록 연결했다.
