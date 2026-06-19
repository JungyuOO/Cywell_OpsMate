# v0.0.49 Phase 1 - OLS Query Payload

## Tasks

- [x] Change CYOps Lightspeed provider request field from `message` to `query`.
- [x] Keep provider/model selection out of the CYOps payload.
- [x] Add a regression test that fails if `message`, `model`, or `provider` are sent to Lightspeed.

## Verification

- `go test ./...` passed.

## Remaining Scope

- Package and install v0.0.49 on CRC.
