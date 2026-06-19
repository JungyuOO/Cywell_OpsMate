# v0.0.45 Phase 1 - Proxy API Path Fix

## Tasks

- [x] Capture `/api/plugins/cyops-console` from the entry script source during load.
- [x] Reuse the captured base from async event handlers.
- [x] Keep direct local appserver fallback paths unchanged.
- [x] Update regression tests.

## Verification

- `go test ./...` passed.
- Appserver tests verify v0.0.45 entry content and proxy-base path logic.

## Remaining Scope

- OLM packaging, image publication, CRC upgrade, and browser chat smoke remain for later phases.
