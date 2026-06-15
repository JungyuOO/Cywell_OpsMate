# Phase 3 - CRD, Sample, And Runbook Update

## Scope

- [x] Add approval field to CRD schema.
- [x] Add approval field to sample CR.
- [x] Update the migration runbook to require approval.
- [x] Link GitHub Issue #77.

## Work Completed

- Updated `config/crd/opsmate.cywell.io_opsmateconfigs.yaml`.
- Updated `config/samples/opsmate_v1alpha1_opsmateconfig.yaml`.
- Updated `spec/0.0.15/pgvector_migration_runbook.md`.

## Verification

- `kubectl kustomize config/crd`

## Remaining Scope

- OpenShift cluster application and runtime status validation remain follow-up work.
