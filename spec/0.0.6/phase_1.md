# Phase 1 - PostgreSQL document repository

## 작업 내용

- [x] `database/sql` 기반 PostgreSQL document repository를 추가한다.
- [x] `List`, `Create`, `Get`, `MarkDeleting`을 구현한다.
- [x] document 원문이 아닌 metadata/object URI만 저장하도록 유지한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] `go test ./...`

## 남은 범위

- [ ] Docker PostgreSQL integration 검증은 Phase 2에서 수행한다.
