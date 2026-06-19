# v0.0.48 Phase 1 - CYOps To Lightspeed Boundary

## Tasks

- [x] Remove LLM provider/model values from the CYOps appserver Lightspeed request payload.
- [x] Stop injecting `CYOPS_LIGHTSPEED_MODEL`, `CYOPS_LIGHTSPEED_PROVIDER`, `LIGHTSPEED_DEFAULT_MODEL`, and `LIGHTSPEED_DEFAULT_PROVIDER` from `OpsMateConfig`.
- [x] Add a regression test that fails if CYOps sends `model` or `provider` to Lightspeed.

## Verification

- `go test ./...` passed.

## Remaining Scope

- Package v0.0.48 and validate on CRC.
