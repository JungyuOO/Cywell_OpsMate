# Phase 3 - runtime docs and operating conditions

## 작업 항목

- [x] migration README를 갱신한다.
- [x] BYTEA fallback과 pgvector activation 조건을 명시한다.
- [x] verification 결과를 문서화한다.

## 검증

- [x] `go test ./...`
- [x] `go build -o .cache\manager.exe ./cmd/manager`

## 남은 범위

- OpenShift cluster에서 pgvector-enabled PostgreSQL image 검증은 후속 버전에서 진행한다.

## 작업 내용

- migration README에 BYTEA fallback, generated `VECTOR(n)` activation path, re-embedding 필요성을 기록했다.
