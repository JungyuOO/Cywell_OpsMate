# v0.0.23 Phase 2 - Diagnostics Schema Hardening

## 작업 내용

- [x] Added `GET /api/ops/diagnostics/schema`.
- [x] Added `ui` metadata to diagnostics response.
- [x] Declared forbidden secret/customer-content fields for future UI consumption.

## 검증

- `go test ./internal/appserver`
- `go test ./...`

## 남은 범위

- v0.0.24 should bind this schema to a visible ConsolePlugin diagnostics view.

## 연결 이슈

- #111
