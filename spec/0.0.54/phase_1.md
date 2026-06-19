# v0.0.54 Phase 1 - Documents Page Handoff

## Tasks

- [x] Remove the injected document workspace overlay from `plugin-entry.js`.
- [x] Change the injected left `자료` action to navigate to `/console-plugin/documents`.
- [x] Keep the CYOps chat drawer focused on chat only.

## Verification

- `go test ./...` passed.
- Plugin entry tests confirm `/console-plugin/documents` navigation and absence of `cyops-doc-workspace`.

## Remaining Scope

- Complete packaging smoke.
