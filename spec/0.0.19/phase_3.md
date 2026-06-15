# Phase 3 - CRD, Sample, And Verification

## Scope

- [x] Add admin token fields to CRD.
- [x] Add runtime status fields to CRD.
- [x] Add admin token Secret fields to sample.
- [x] Run full test/build verification.
- [x] Link GitHub Issue #92.

## Work Completed

- Updated `config/crd/opsmate.cywell.io_opsmateconfigs.yaml`.
- Updated `config/samples/opsmate_v1alpha1_opsmateconfig.yaml`.
- Verified CRD with kustomize.

## Verification

- `kubectl kustomize config/crd`
- `go test ./...`
- `go build -o .cache\manager.exe ./cmd/manager`
- `go build -o .cache\cyops-pgvector-migrate.exe ./cmd/cyops-pgvector-migrate`

## Remaining Scope

- Apply CRD/sample to target OpenShift cluster.
