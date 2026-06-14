# Phase 3 - ingestion service flow

## 작업 항목

- [x] document id 기준 ingestion orchestration을 추가한다.
- [x] 원본 파일 읽기, parsing, chunking, persistence를 연결한다.
- [x] 실패 시 document status와 last_error를 갱신한다.

## 검증

- [x] `go test ./internal/appserver`
- [x] Docker PostgreSQL integration test

## 남은 범위

- async queue와 별도 worker Deployment는 후속 버전에서 검토한다.

## 작업 내용

- `IngestionService.IngestDocument`가 document 조회, processing 전환, parsing, chunk 저장, ready/failed 전환을 수행한다.
- Docker PostgreSQL integration test로 local storage 원본 파일에서 chunk table 저장까지 검증했다.
