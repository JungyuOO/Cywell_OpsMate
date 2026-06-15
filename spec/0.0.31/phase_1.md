# v0.0.31 Phase 1 - Isolated Local OLM Manifests

## Work

- [x] Added `deploy/olm/local-crc`.
- [x] Used namespace `cywell-opsmate-olm` to avoid direct-bootstrap conflicts.
- [x] Added local OperatorGroup and Subscription manifests.
- [x] Added CRC-only `OpsMateConfig` and PostgreSQL password Secret for reconcile smoke.

## Verification

- `kubectl kustomize deploy/olm/local-crc`
- `oc apply --dry-run=client -k deploy/olm/local-crc`

## Remaining Scope

- None for manifest creation.

## Linked Issue

- #150
