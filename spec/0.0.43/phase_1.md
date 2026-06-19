# v0.0.43 Phase 1 - Callback Module Map

## Tasks

- [x] Add a supported `console.flag` extension for CYOps launcher execution.
- [x] Register callback plugin entry as `cyops-console@0.0.43`.
- [x] Expose `cyopsLauncherFlag` through a callback module map.
- [x] Update regression tests.

## Verification

- `go test ./...` passed.
- Appserver tests verify `console.flag`, `$codeRef: "cyopsLauncherFlag"`, `cyops-console@0.0.43`, and `cyopsLauncherFlag` in the served plugin content.

## Remaining Scope

- OLM packaging and CRC smoke remain for Phase 2.
