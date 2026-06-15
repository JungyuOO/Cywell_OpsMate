# v0.0.32 Phase 2 - OLM Upgrade

## Work

- [x] Confirmed Subscription saw `cywell-opsmate.v0.0.31`.
- [x] Reviewed and approved generated upgrade InstallPlan `install-2rf64`.
- [x] Confirmed CSV `cywell-opsmate.v0.0.31` reached `Succeeded` after a local recovery patch moved the OLM-managed deployment image to v0.0.31.

## Verification

- `oc get subscription,installplan,csv,deploy,pod -n cywell-opsmate-olm`
- `oc get installplan -n cywell-opsmate-olm -o yaml`
- `oc patch installplan install-2rf64 -n cywell-opsmate-olm --type=merge --patch '{\"spec\":{\"approved\":true}}'`
- `oc set image deployment/cywell-opsmate-controller-manager manager=ghcr.io/jungyuoo/cywell-opsmate:v0.0.31 -n cywell-opsmate-olm`

Observed evidence:

- InstallPlan `install-2rf64` referenced `cywell-opsmate.v0.0.31`, bundle `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.31`, and `replaces: cywell-opsmate.v0.0.29`.
- OLM initially held v0.0.31 at `Pending/RequirementsUnknown` because v0.0.29 was unhealthy.
- After the local recovery image patch, CSV `cywell-opsmate.v0.0.31` reached `Succeeded` and deployment image was `ghcr.io/jungyuoo/cywell-opsmate:v0.0.31`.

## Remaining Scope

- Fix the namespace cache/RBAC mismatch observed in upgraded manager logs.

## Linked Issue

- #158
