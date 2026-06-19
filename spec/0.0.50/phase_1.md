# v0.0.50 Phase 1 - Strict OLS Request Body

## Tasks

- [x] Remove `context` and `clusterContext` from Lightspeed provider payload.
- [x] Keep `query` as the only OLS request body field.
- [x] Add regression assertions for rejected extra fields.

## Verification

- `go test ./...` passed.

## Remaining Scope

- Package, upgrade, and smoke on CRC.
