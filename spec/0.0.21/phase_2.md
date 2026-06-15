# v0.0.21 Phase 2 - Live Migration Job Reconciliation

## 작업 내용

- [x] Reconciled `PGVectorMigrationJob` when `pgVectorMigrationApproved=true`.
- [x] Added Job reconciliation support to the generic object reconciler.
- [x] Applied live Job status back to `OpsMateConfig.status`.
- [x] Recorded builder errors in `status.pgVectorLastError`.

## 검증

- `go test ./internal/controller`
- `go test ./internal/controller/postgres`

## 남은 범위

- v0.0.22 should add a cluster smoke command or script that waits for Job completion and fetches sanitized logs.

## 연결 이슈

- #101
