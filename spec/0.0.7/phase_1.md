# Phase 1 - appserver runtime config

## 작업 항목

- [x] 환경변수 기반 appserver config loader를 추가한다.
- [x] PostgreSQL DSN이 있을 때 repository를 구성하는 builder를 추가한다.
- [x] document storage path와 Lightspeed endpoint를 dependency에 연결한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- upload handler에서 storage와 repository를 하나의 흐름으로 연결했다.

## 작업 내용

- `LoadConfigFromEnv`와 `NewServerFromConfig`를 추가했다.
- PostgreSQL DSN이 있을 때 migration, repository, storage, Lightspeed provider를 runtime dependency로 구성한다.
