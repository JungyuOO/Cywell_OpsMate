# v0.0.44 Phase 1 - Entry Load Mount

## Tasks

- [x] Keep `console.flag` and callback module map registration.
- [x] Register callback plugin entry as `cyops-console@0.0.44`.
- [x] Call the CYOps launcher mount path immediately when `/plugin-entry.js` is executed.
- [x] Add `data-cyops-plugin-entry="0.0.44"` as a DOM execution marker.
- [x] Add delayed mount retries for Console page re-render timing.
- [x] Update regression tests.

## Verification

- `go test ./...` passed.
- Appserver tests verify the v0.0.44 entry string, plugin registration, execution marker, direct mount retry, and launcher content.

## Remaining Scope

- Image publication, CRC upgrade, and authenticated browser confirmation remain for later phases.
