# v0.0.33 Phase 2 - OLM Upgrade To v0.0.32

## Work

- [x] Reviewed generated v0.0.32 InstallPlan `install-thxvg`.
- [x] Approved v0.0.32 InstallPlan.
- [x] Confirmed CSV and deployment reached v0.0.32 after a local recovery image patch from the known-crashing v0.0.31 deployment.

## Verification

- `oc get subscription,installplan,csv,deploy,pod -n cywell-opsmate-olm`
- `oc get installplan -n cywell-opsmate-olm -o yaml`
- `oc patch installplan install-thxvg -n cywell-opsmate-olm --type=merge --patch '{\"spec\":{\"approved\":true}}'`
- `oc set image deployment/cywell-opsmate-controller-manager manager=ghcr.io/jungyuoo/cywell-opsmate:v0.0.32 -n cywell-opsmate-olm`

Observed evidence:

- InstallPlan `install-thxvg` referenced `cywell-opsmate.v0.0.32`, bundle `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.32`, and `replaces: cywell-opsmate.v0.0.31`.
- CSV `cywell-opsmate.v0.0.32` reached `Succeeded`.
- Deployment image was `ghcr.io/jungyuoo/cywell-opsmate:v0.0.32` with `POD_NAMESPACE` env.

## Remaining Scope

- Retry local `OpsMateConfig` reconcile and capture remaining blockers.

## Linked Issue

- #162
