# Phase 3 - CRD, Sample, And Runbook

## Scope

- [x] Add database and embedding pgvector fields to the CRD schema.
- [x] Update the sample `OpsMateConfig`.
- [x] Write an OpenShift pgvector migration and smoke runbook.
- [x] Link GitHub Issue #72.

## Work Completed

- Updated `config/crd/opsmate.cywell.io_opsmateconfigs.yaml`.
- Updated `config/samples/opsmate_v1alpha1_opsmateconfig.yaml`.
- Added `pgvector_migration_runbook.md`.

## Verification

- `go test ./...`
- `kubectl kustomize config/crd`

## Remaining Scope

- Run the runbook against the target OpenShift cluster.
