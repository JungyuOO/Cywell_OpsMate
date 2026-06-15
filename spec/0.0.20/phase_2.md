# v0.0.20 Phase 2 - Migration Job Runtime Evidence

## 작업 내용

- [x] Added `ApplyPGVectorMigrationJobStatus`.
- [x] Marked `status.pgVectorReady=true` when `JobComplete=True`.
- [x] Marked `status.pgVectorReady=false` and recorded `status.pgVectorLastError` when `JobFailed=True`.

## 검증

- `go test ./internal/controller/postgres`

## 남은 범위

- v0.0.21 should wire Job lookup/reconcile execution against an OpenShift cluster instead of only providing the evidence mapper.

## 연결 이슈

- #96
