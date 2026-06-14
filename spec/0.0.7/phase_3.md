# Phase 3 - runtime migration strategy

## 작업 항목

- [x] migration SQL을 appserver runtime에서 적용할 수 있는 runner를 추가한다.
- [x] migration runner가 반복 실행 가능한지 테스트한다.
- [x] integration test에서 migration runner를 통해 schema를 준비한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] Docker PostgreSQL integration test

## 남은 범위

- schema version table과 online migration orchestration은 운영 요구가 명확해진 뒤 보강한다.

## 작업 내용

- embedded SQL 기반 `ApplyMigrations`를 추가했다.
- Docker PostgreSQL에서 migration 반복 적용 후 upload persistence flow를 검증했다.
