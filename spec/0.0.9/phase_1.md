# Phase 1 - embedding provider boundary

## 작업 항목

- [x] embedding provider interface를 추가한다.
- [x] deterministic mock provider를 구현한다.
- [x] provider output이 안정적인지 테스트한다.

## 검증

- [x] `go test ./internal/appserver`

## 남은 범위

- embedding persistence는 Phase 2에서 구현했다.

## 작업 내용

- `EmbeddingProvider`와 `DeterministicEmbeddingProvider`를 추가했다.
- mock provider는 SHA-256 기반 deterministic byte vector를 생성한다.
