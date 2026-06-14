# Phase 3 - chat retrieval API

## 작업 항목

- [x] retriever interface를 추가한다.
- [x] `/api/chat` RAG 요청에서 retriever를 호출한다.
- [x] citations를 provider context와 chat response에 연결한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- vector similarity ranking은 후속 버전에서 구현한다.

## 작업 내용

- `Retriever`와 `PostgresRetriever`를 추가했다.
- `/api/chat`의 `rag.enabled=true` path가 retriever context를 provider request에 붙이고 citations를 응답한다.
