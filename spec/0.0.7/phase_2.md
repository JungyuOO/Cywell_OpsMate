# Phase 2 - upload storage and metadata persistence

## 작업 항목

- [x] `POST /api/documents`가 업로드 파일을 `DocumentStorage`에 저장한다.
- [x] 저장된 object URI와 size를 repository metadata에 기록한다.
- [x] 저장 실패 시 metadata가 생성되지 않도록 handler error path를 검증한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] PostgreSQL DSN 기반 integration test

## 남은 범위

- chunk/parser/embedding worker 구현은 v0.0.8로 넘긴다.

## 작업 내용

- 저장소가 설정된 업로드 경로에서 파일 저장 후 metadata를 생성하도록 handler를 확장했다.
- PostgreSQL repository가 `object_uri`를 insert/list/get/delete marker 응답에 포함하도록 확장했다.
