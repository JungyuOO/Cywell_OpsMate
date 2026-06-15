# v0.0.24 Phase 1 - Console Diagnostics Served View

## 작업 내용

- [x] Added `/console-plugin/diagnostics`.
- [x] Added `/console-plugin/diagnostics.js`.
- [x] Added `/console-plugin/diagnostics.css`.
- [x] Kept the view inside the appserver/ConsolePlugin backend path.

## 검증

- `go test ./internal/appserver`

## 남은 범위

- Future versions can replace the static view with a full bundled dynamic plugin UI.

## 연결 이슈

- #115
