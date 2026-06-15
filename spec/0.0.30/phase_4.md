# v0.0.30 Phase 4 - RBAC Correction And Local Handoff

## Work

- [x] Moved direct-bootstrap ConsolePlugin permissions to ClusterRole/ClusterRoleBinding.
- [x] Moved CSV ConsolePlugin permissions to `clusterPermissions`.
- [x] Captured next local OLM install steps.

## Verification

- `kubectl kustomize config/default`
- `oc apply --dry-run=client -f deploy/olm/bundle/manifests`
- `go test ./...`

## Remaining Scope

- See `v0.0.31_handoff.md`.

## Linked Issue

- #148
