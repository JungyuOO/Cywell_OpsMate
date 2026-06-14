# Phase 1 - pgvector readiness and config

## 작업 항목

- [x] `CYOPS_PGVECTOR_REQUIRED` config를 추가한다.
- [x] pgvector readiness validation을 추가한다.
- [x] readiness failure가 명확한 error로 드러나게 한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] Docker PostgreSQL integration evidence

## 남은 범위

- 실제 `VECTOR(n)` migration은 후속 버전에서 진행한다.

## 작업 내용

- `CYOPS_PGVECTOR_REQUIRED`와 embedding provider config env를 추가했다.
- `CheckPGVectorReady`가 extension 준비 실패를 `pgvector extension is not ready` error로 감싼다.
