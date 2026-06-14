# Phase 2 - pgvector migration helper

## 작업 항목

- [x] embedding dimension validation을 추가한다.
- [x] `VECTOR(n)` 전환 SQL generator를 추가한다.
- [x] SQL helper tests를 추가한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- migration 자동 적용은 v0.0.12 이후 runtime activation에서 검토한다.

## 작업 내용

- `PGVectorEmbeddingMigrationSQL`을 추가했다.
- dimensions validation과 `VECTOR(n)` SQL generation을 테스트했다.
