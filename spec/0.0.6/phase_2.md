# Phase 2 - Docker PostgreSQL integration verification

## 작업 내용

- [x] Docker PostgreSQL 컨테이너를 준비한다.
- [x] migration SQL을 적용한다.
- [x] repository create/list/get/delete marker를 실제 PostgreSQL에서 검증한다.

## 검증

- [x] `CYOPS_POSTGRES_TEST_DSN=postgres://cyops:cyops@localhost:55432/cyops?sslmode=disable go test ./internal/appserver -run TestPostgresDocumentRepositoryIntegration -count=1 -v`
- [x] `go test ./...`

## 남은 범위

- [ ] OpenShift PostgreSQL 운영화는 후속 버전으로 이관한다.
