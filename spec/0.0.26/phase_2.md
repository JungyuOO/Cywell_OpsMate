# v0.0.26 Phase 2 - Deployment Command And Listen Wiring

## Work

- [x] Set the appserver container command to `/appserver`.
- [x] Added `CYOPS_LISTEN_ADDRESS=:8443` to the appserver Deployment env.
- [x] Kept existing service-ca TLS mount and env paths.

## Verification

- `go test ./internal/controller/appserver`
- `kubectl kustomize config/crd`

## Remaining Scope

- Apply the reconciled Deployment to a real OpenShift cluster in the next live evidence pass.

## Linked Issue

- #126
