# v0.0.22 Phase 2 - Migration Job Operational Smoke

## 작업 내용

- [x] Added optional wait for `<name>-pgvector-migration`.
- [x] Added migration Job log readback.
- [x] Added sensitive output pattern checks.
- [x] Added `status.pgVectorReady=true` assertion after Job completion.

## 검증

- PowerShell parser validation for `deploy/scripts/openshift-v022-smoke.ps1`.
- Existing controller tests cover Job reconciliation/status mapping.

## 남은 범위

- Run the script against a real OpenShift cluster once image and DSN Secret values are available.

## 연결 이슈

- #106
