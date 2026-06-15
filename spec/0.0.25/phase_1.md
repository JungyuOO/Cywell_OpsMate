# v0.0.25 Phase 1 - Local Appserver Executable

## Work

- [x] Added `cmd/appserver`.
- [x] Added `CYOPS_LISTEN_ADDRESS` support.
- [x] Kept memory-backed mode usable without PostgreSQL.

## Verification

- `go test ./cmd/appserver`
- `go build -o .cache\cyops-appserver.exe ./cmd/appserver`

## Remaining Scope

- Production image entrypoint wiring remains a release packaging follow-up.

## Linked Issue

- #120
