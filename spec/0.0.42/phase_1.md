# v0.0.42 Phase 1 - Gateway Resources

## Tasks

- [x] Add gateway controller resource builders.
- [x] Reconcile gateway ConfigMap, Deployment, and Service.
- [x] Point ConsolePlugin backend to the gateway Service.
- [x] Add regression tests for gateway resources and ConsolePlugin backend service.

## Verification

- `go test ./...` passed.
- Gateway unit tests verify nginx TLS listener, appserver proxy target, service-serving cert annotation, and gateway labels.
- ConsolePlugin unit tests verify backend Service name changes to `sample-gateway`.

## Remaining Scope

- OLM publish and CRC smoke remain for Phase 2.
