# v0.0.22 Phase 3 - Admin Diagnostics Endpoint

## 작업 내용

- [x] Added `GET /api/ops/diagnostics`.
- [x] Protected diagnostics with the same admin authorization as re-embedding.
- [x] Returned retrieval metrics, document aggregate counts, admin identity headers, re-embedding availability, and endpoint links.
- [x] Avoided prompts, document content, DSNs, and tokens.

## 검증

- `go test ./internal/appserver`
- `go test ./...`

## 남은 범위

- v0.0.23 should add a ConsolePlugin admin diagnostics view that consumes this endpoint.

## 연결 이슈

- #107
