# Phase 1 - EmbeddingSpec API and env boundary

## 작업 항목

- [x] `EmbeddingSpec`을 API type에 추가한다.
- [x] appserver Deployment env에 embedding config를 전달한다.
- [x] SecretKeyRef로 embedding token boundary를 추가한다.

## 검증

- [x] `go test ./internal/controller/appserver`
- [x] `go test ./api/v1alpha1`

## 남은 범위

- pgvector migration helper는 Phase 2에서 구현했다.

## 작업 내용

- `EmbeddingSpec`을 추가하고 appserver Deployment env/SecretKeyRef에 연결했다.
- appserver runtime이 `CYOPS_EMBEDDING_TOKEN`을 읽고 HTTP embedding provider Authorization header에만 사용한다.
