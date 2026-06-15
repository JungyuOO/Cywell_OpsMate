# v0.0.29 Phase 1 - OLM Bundle Manifests

## Work

- [x] Added CYOps ClusterServiceVersion.
- [x] Included `OpsMateConfig` CRD in the bundle manifests.
- [x] Declared owned CRD, install modes, permissions, related images, and example CR.

## Verification

- `oc apply --dry-run=client -f deploy/olm/bundle/manifests`

## Remaining Scope

- Validate the bundle with operator tooling once `operator-sdk` or `opm` is part of CI.

## Linked Issue

- #140
