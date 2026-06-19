# v0.0.51 Phase 1 - Scoped Internal Service TLS

## Tasks

- [x] Detect OpenShift internal service DNS endpoints.
- [x] Use a dedicated HTTP transport for HTTPS `.svc` endpoints.
- [x] Keep external endpoints on the default Go HTTP client.
- [x] Add regression tests for scoped behavior.

## Verification

- `go test ./...` passed.

## Remaining Scope

- Package v0.0.51 and run CRC smoke.
