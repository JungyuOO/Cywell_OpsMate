# v0.0.26 Phase 1 - Appserver TLS Runtime

## Work

- [x] Added `TLS_CERT_FILE` and `TLS_KEY_FILE` handling to `cmd/appserver`.
- [x] Kept local development on plain HTTP when TLS files are not configured.
- [x] Added startup validation that rejects partial TLS configuration.

## Verification

- `go test ./cmd/appserver`
- `go build -o .cache\cyops-appserver.exe ./cmd/appserver`

## Remaining Scope

- Live service-ca certificate verification requires an OpenShift deployment.

## Linked Issue

- #125
