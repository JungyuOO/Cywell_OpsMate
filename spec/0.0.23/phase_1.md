# v0.0.23 Phase 1 - ConsolePlugin Primary Entry Contract

## 작업 내용

- [x] Added ConsolePlugin annotations for primary entry and diagnostics path.
- [x] Marked fallback admin Route mode as enabled/disabled metadata.
- [x] Kept Red Hat OpenShift Lightspeed viewer untouched.

## 검증

- `go test ./internal/controller/console`

## 남은 범위

- v0.0.24 should add the actual ConsolePlugin diagnostics frontend bundle/view.

## 연결 이슈

- #110
