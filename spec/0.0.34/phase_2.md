# v0.0.34 Phase 2 - OLM Upgrade To v0.0.33

## Work

- [x] Reviewed generated v0.0.33 InstallPlan `install-5r22q`.
- [x] Approved v0.0.33 InstallPlan.
- [x] Confirmed CSV and deployment are v0.0.33.

## Verification

- `oc get subscription,installplan,csv,deploy,pod -n cywell-opsmate-olm`
- `oc get installplan -n cywell-opsmate-olm -o yaml`
- `oc patch installplan install-5r22q -n cywell-opsmate-olm --type=merge --patch '{\"spec\":{\"approved\":true}}'`

Observed evidence:

- InstallPlan `install-5r22q` referenced `cywell-opsmate.v0.0.33`, bundle `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.33`, and `replaces: cywell-opsmate.v0.0.32`.
- CSV `cywell-opsmate.v0.0.33` reached `Succeeded`.
- Manager deployment image was `ghcr.io/jungyuoo/cywell-opsmate:v0.0.33`.

## Remaining Scope

- Retry local `OpsMateConfig` reconcile and capture remaining blockers.

## Linked Issue

- #168
